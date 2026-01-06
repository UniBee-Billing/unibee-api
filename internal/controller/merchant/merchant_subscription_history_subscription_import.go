package merchant

import (
	"context"
	"fmt"
	"strings"
	"unibee/api/merchant/subscription"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/plan/period"
	detailService "unibee/internal/logic/subscription/service/detail"
	user2 "unibee/internal/logic/user"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func (c *ControllerSubscription) HistorySubscriptionImport(ctx context.Context, req *subscription.HistorySubscriptionImportReq) (res *subscription.HistorySubscriptionImportRes, err error) {
	// Get merchant context
	merchantId := _interface.GetMerchantId(ctx)

	// Validate required fields
	utility.Assert(len(req.ExternalSubscriptionId) > 0, "ExternalSubscriptionId is blank")
	utility.Assert(len(req.ExternalPlanId) > 0 || req.PlanId > 0, "one of planId or ExternalPlanId is required")
	utility.Assert(len(req.Gateway) > 0, "Gateway is blank")
	utility.Assert(len(req.CurrentPeriodStart) > 0, "CurrentPeriodStart is blank")
	utility.Assert(len(req.CurrentPeriodEnd) > 0, "CurrentPeriodEnd is blank")
	utility.Assert(req.TotalAmount > 0, "TotalAmount should greater than 0")

	// Validate countryCode if provided
	if len(req.CountryCode) > 0 {
		err = utility.ValidateCountryCode(req.CountryCode)
		utility.AssertError(err, "Invalid country code: "+req.CountryCode)
	}

	// Find or create user
	user, err := user2.QueryOrCreateUser(ctx, &user2.NewUserInternalReq{
		ExternalUserId: req.ExternalUserId,
		Email:          req.Email,
		CountryCode:    req.CountryCode,
		MerchantId:     _interface.GetMerchantId(ctx),
	})
	utility.AssertError(err, "QueryOrCreateUser failed")
	utility.Assert(user != nil, "can't find user by ExternalUserId or Email")

	// Get tax percentage from request instead of user
	taxPercentage := req.TaxPercentage
	if !vat_gateway.GetDefaultVatGateway(ctx, user.MerchantId).VatRatesEnabled() && req.TaxPercentage > 0 {
		taxPercentage = req.TaxPercentage
	} else if vat_gateway.GetDefaultVatGateway(ctx, user.MerchantId).VatRatesEnabled() {
		// Use user's tax percentage if VAT gateway is enabled and no specific tax percentage provided
		taxPercentage = user.TaxPercentage
	}

	// Find plan by external plan ID
	var plan *entity.Plan
	if req.PlanId > 0 {
		plan = query.GetPlanById(ctx, req.PlanId)
	} else if len(req.ExternalPlanId) > 0 {
		plan = query.GetPlanByExternalPlanId(ctx, merchantId, req.ExternalPlanId)
	}
	utility.Assert(plan != nil, "can't find plan by ExternalPlanId or planId")
	utility.Assert(plan.Status != consts.PlanStatusEditable, "plan status should not editable")

	// Validate gateway
	gatewayImpl := api.GatewayNameMapping[req.Gateway]
	utility.Assert(gatewayImpl != nil, "Invalid Gateway, should be one of "+strings.Join(api.ExportGatewaySetupListKeys(), "|"))
	gateway := query.GetDefaultGatewayByGatewayName(ctx, merchantId, req.Gateway)
	utility.Assert(gateway != nil, fmt.Sprintf("Gateway:%s not found, please setup first", req.Gateway))
	gatewayId := gateway.Id

	// Parse quantity
	if req.Quantity == 0 {
		req.Quantity = 1
	}

	// Parse time fields
	currentPeriodStart := gtime.New(req.CurrentPeriodStart)
	currentPeriodEnd := gtime.New(req.CurrentPeriodEnd)

	// Validate time relationships for history subscription
	now := gtime.Now()
	utility.Assert(currentPeriodStart.Timestamp() < now.Timestamp(), "CurrentPeriodStart should earlier then now")
	utility.Assert(currentPeriodEnd.Timestamp() < now.Timestamp(), "CurrentPeriodEnd should earlier then now")
	utility.Assert(currentPeriodEnd.Timestamp() > currentPeriodStart.Timestamp(), "currentPeriodEnd should later then currentPeriodStart")

	// Check if subscription already exists
	existingSubscription := query.GetSubscriptionByExternalSubscriptionId(ctx, req.ExternalSubscriptionId)
	tag := ""

	if _interface.Context().Get(ctx).MerchantMember != nil {
		tag = fmt.Sprintf("ImportByMember:%d", _interface.Context().Get(ctx).MerchantMember.Id)
	} else {
		tag = fmt.Sprintf("ImportByOpenAPI")
	}
	if existingSubscription != nil {
		// Check permission to override
		utility.Assert(existingSubscription.Data == tag, "no permission to override, "+existingSubscription.Data)
		utility.Assert(existingSubscription.UserId == user.Id, "no permission to override, user not match")
		var dunningTime = period.GetDunningTimeFromEnd(ctx, utility.MaxInt64(currentPeriodEnd.Timestamp(), 0), plan.Id)
		// Update existing subscription
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().Type:                        consts.SubTypeUniBeeControl,
			dao.Subscription.Columns().Status:                      consts.SubStatusExpired,
			dao.Subscription.Columns().Amount:                      req.TotalAmount,
			dao.Subscription.Columns().Currency:                    plan.Currency,
			dao.Subscription.Columns().PlanId:                      plan.Id,
			dao.Subscription.Columns().Quantity:                    req.Quantity,
			dao.Subscription.Columns().GatewayId:                   gatewayId,
			dao.Subscription.Columns().GatewayItemData:             "",
			dao.Subscription.Columns().GatewayDefaultPaymentMethod: "",
			dao.Subscription.Columns().TaxPercentage:               taxPercentage,
			dao.Subscription.Columns().BillingCycleAnchor:          currentPeriodStart.Timestamp(),
			dao.Subscription.Columns().CurrentPeriodStart:          currentPeriodStart.Timestamp(),
			dao.Subscription.Columns().CurrentPeriodEnd:            currentPeriodEnd.Timestamp(),
			dao.Subscription.Columns().DunningTime:                 dunningTime,
			dao.Subscription.Columns().CurrentPeriodStartTime:      currentPeriodStart,
			dao.Subscription.Columns().CurrentPeriodEndTime:        currentPeriodEnd,
			dao.Subscription.Columns().FirstPaidTime:               currentPeriodStart.Timestamp(),
			dao.Subscription.Columns().CreateTime:                  currentPeriodStart.Timestamp(),
		}).Where(dao.Subscription.Columns().Id, existingSubscription.Id).OmitNil().Update()
		utility.AssertError(err, "Update subscription error")

		// Log operation
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     existingSubscription.MerchantId,
			Target:         fmt.Sprintf("SubscriptionHistory(%s)", existingSubscription.SubscriptionId),
			Content:        "ImportOverride",
			UserId:         existingSubscription.UserId,
			SubscriptionId: existingSubscription.SubscriptionId,
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)

		// Get updated subscription detail
		subscriptionDetail, err := detailService.SubscriptionDetail(ctx, existingSubscription.SubscriptionId)
		utility.AssertError(err, "Get subscription detail error")

		return &subscription.HistorySubscriptionImportRes{
			Subscription: subscriptionDetail,
		}, nil
	} else {
		var dunningTime = period.GetDunningTimeFromEnd(ctx, utility.MaxInt64(currentPeriodEnd.Timestamp(), 0), plan.Id)
		// Create new history subscription
		newSubscription := &entity.Subscription{
			Type:                        consts.SubTypeUniBeeControl,
			SubscriptionId:              utility.CreateSubscriptionId(),
			ExternalSubscriptionId:      req.ExternalSubscriptionId,
			UserId:                      user.Id,
			Amount:                      req.TotalAmount,
			Currency:                    plan.Currency,
			MerchantId:                  merchantId,
			PlanId:                      plan.Id,
			Quantity:                    req.Quantity,
			GatewayId:                   gatewayId,
			Status:                      consts.SubStatusExpired,
			CurrentPeriodStart:          currentPeriodStart.Timestamp(),
			CurrentPeriodEnd:            currentPeriodEnd.Timestamp(),
			CurrentPeriodStartTime:      currentPeriodStart,
			CurrentPeriodEndTime:        currentPeriodEnd,
			DunningTime:                 dunningTime,
			BillingCycleAnchor:          currentPeriodStart.Timestamp(),
			FirstPaidTime:               currentPeriodStart.Timestamp(),
			CreateTime:                  currentPeriodStart.Timestamp(),
			CountryCode:                 user.CountryCode,
			VatNumber:                   user.VATNumber,
			TaxPercentage:               taxPercentage,
			GatewaySubscriptionId:       req.ExternalSubscriptionId,
			GatewayItemData:             "",
			Data:                        tag,
			CurrentPeriodPaid:           1,
			GatewayDefaultPaymentMethod: "",
		}

		result, err := dao.Subscription.Ctx(ctx).Data(newSubscription).OmitNil().Insert(newSubscription)
		utility.AssertError(err, "Save subscription error")
		id, err := result.LastInsertId()
		utility.AssertError(err, "Get last insert id error")
		newSubscription.Id = uint64(id)

		// Log operation
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     newSubscription.MerchantId,
			Target:         fmt.Sprintf("SubscriptionHistory(%s)", newSubscription.SubscriptionId),
			Content:        "ImportNew",
			UserId:         newSubscription.UserId,
			SubscriptionId: newSubscription.SubscriptionId,
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)

		// Create subscription timeline
		uniqueId := fmt.Sprintf("%s-%v-%v-%v", newSubscription.ExternalSubscriptionId, newSubscription.PlanId, newSubscription.CurrentPeriodStart, newSubscription.CurrentPeriodEnd)
		timeline := query.GetSubscriptionTimeLineByUniqueId(ctx, uniqueId)
		utility.Assert(timeline == nil, "same history record exist: "+uniqueId)

		// Create timeline record
		timeline = &entity.SubscriptionTimeline{
			MerchantId:      newSubscription.MerchantId,
			UserId:          newSubscription.UserId,
			SubscriptionId:  newSubscription.SubscriptionId,
			UniqueId:        uniqueId,
			Currency:        newSubscription.Currency,
			PlanId:          newSubscription.PlanId,
			Quantity:        newSubscription.Quantity,
			AddonData:       newSubscription.AddonData,
			Status:          consts.SubTimeLineStatusFinished,
			GatewayId:       newSubscription.GatewayId,
			PeriodStart:     newSubscription.CurrentPeriodStart,
			PeriodEnd:       newSubscription.CurrentPeriodEnd,
			PeriodStartTime: gtime.NewFromTimeStamp(newSubscription.CurrentPeriodStart),
			PeriodEndTime:   gtime.NewFromTimeStamp(newSubscription.CurrentPeriodEnd),
			CreateTime:      gtime.Now().Timestamp(),
		}

		timelineResult, err := dao.SubscriptionTimeline.Ctx(ctx).Data(timeline).OmitNil().Insert(timeline)
		utility.AssertError(err, "Save timeline error")
		timelineId, err := timelineResult.LastInsertId()
		utility.AssertError(err, "Get timeline last insert id error")
		timeline.Id = uint64(timelineId)

		// Get subscription detail
		subscriptionDetail, err := detailService.SubscriptionDetail(ctx, newSubscription.SubscriptionId)
		utility.AssertError(err, "Get subscription detail error")

		return &subscription.HistorySubscriptionImportRes{
			Subscription: subscriptionDetail,
		}, nil
	}
}
