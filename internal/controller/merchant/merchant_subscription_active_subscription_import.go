package merchant

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"
	"unibee/api/bean"
	"unibee/api/merchant/subscription"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/gateway/api"
	handler2 "unibee/internal/logic/invoice/handler"
	service3 "unibee/internal/logic/invoice/service"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/plan/period"
	detailService "unibee/internal/logic/subscription/service/detail"
	"unibee/internal/logic/subscription/timeline"
	user2 "unibee/internal/logic/user"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func (c *ControllerSubscription) ActiveSubscriptionImport(ctx context.Context, req *subscription.ActiveSubscriptionImportReq) (res *subscription.ActiveSubscriptionImportRes, err error) {
	// Get merchant context
	merchantId := _interface.GetMerchantId(ctx)

	// Validate required fields
	utility.Assert(len(req.ExternalSubscriptionId) > 0, "ExternalSubscriptionId is blank")
	utility.Assert(len(req.ExternalPlanId) > 0 || req.PlanId > 0, "one of planId or ExternalPlanId is required")
	utility.Assert(len(req.Gateway) > 0, "Gateway is blank")
	utility.Assert(len(req.CurrentPeriodStart) > 0, "CurrentPeriodStart is blank")
	utility.Assert(len(req.CurrentPeriodEnd) > 0, "CurrentPeriodEnd is blank")
	utility.Assert(len(req.BillingCycleAnchor) > 0, "BillingCycleAnchor is blank")
	utility.Assert(len(req.CreateTime) > 0, "CreateTime is blank")

	// Validate countryCode if provided
	if len(req.CountryCode) > 0 {
		err = utility.ValidateCountryCode(req.CountryCode)
		utility.AssertError(err, "Invalid country code: "+req.CountryCode)
	}

	// Find user by external user ID
	user, err := user2.QueryOrCreateUser(ctx, &user2.NewUserInternalReq{
		ExternalUserId: req.ExternalUserId,
		Email:          req.Email,
		CountryCode:    req.CountryCode,
		VATNumber:      req.VatNumber,
		MerchantId:     _interface.GetMerchantId(ctx),
	})
	utility.AssertError(err, "QueryOrCreateUser failed")
	utility.Assert(user != nil, "can't find user by ExternalUserId or Email")

	// Find plan by external plan ID
	var plan *entity.Plan
	if req.PlanId > 0 {
		plan = query.GetPlanById(ctx, req.PlanId)
	} else if len(req.ExternalPlanId) > 0 {
		plan = query.GetPlanByExternalPlanId(ctx, merchantId, req.ExternalPlanId)
	}
	utility.Assert(plan != nil, "can't find plan by ExternalPlanId or planId")
	utility.Assert(plan.Status != consts.PlanStatusEditable, "plan status should not editable")
	anotherSub := query.GetLatestActiveOrIncompleteSubscriptionByUserId(ctx, user.Id, merchantId, plan.ProductId)
	utility.Assert(anotherSub == nil || anotherSub.ExternalSubscriptionId == req.ExternalSubscriptionId, "User Already have anther active subscription in same product")

	// Validate gateway
	gatewayImpl := api.GatewayNameMapping[req.Gateway]
	utility.Assert(gatewayImpl != nil, "Invalid Gateway, should be one of "+strings.Join(api.ExportGatewaySetupListKeys(), "|"))
	gateway := query.GetDefaultGatewayByGatewayName(ctx, merchantId, req.Gateway)
	utility.Assert(gateway != nil, fmt.Sprintf("Gateway:%s not found, please setup first", req.Gateway))
	gatewayId := gateway.Id

	if req.Quantity == 0 {
		req.Quantity = 1
	}

	countryCode := user.CountryCode
	vatNumber := user.VATNumber
	if req.CountryCode != "" {
		countryCode = req.CountryCode
	}
	if req.VatNumber != "" {
		vatNumber = req.VatNumber
	}
	taxPercentage := user.TaxPercentage
	if vat_gateway.GetDefaultVatGateway(ctx, user.MerchantId).VatRatesEnabled() {
		taxPercentage, _ = vat_gateway.ComputeMerchantVatPercentage(ctx, merchantId, countryCode, gatewayId, vatNumber)
	} else {
		taxPercentage = req.TaxPercentage
	}

	totalAmountExcludingTax := plan.Amount * req.Quantity
	var taxAmount = int64(math.Round(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(taxPercentage)))
	totalAmount := totalAmountExcludingTax + taxAmount
	if req.ExpectedTotalAmount > 0 && req.ExpectedTotalAmount != totalAmount {
		utility.Assert(false, fmt.Sprintf("ExpectedTotalAmount Verify failed, calculated totalAmount %s != expectedTotalAmount %s", utility.ConvertCentToDollarStr(totalAmount, plan.Currency), utility.ConvertCentToDollarStr(req.ExpectedTotalAmount, plan.Currency)))
	}

	// Parse time fields
	currentPeriodStart := gtime.New(req.CurrentPeriodStart)
	currentPeriodEnd := gtime.New(req.CurrentPeriodEnd)
	billingCycleAnchor := gtime.New(req.BillingCycleAnchor)
	createTime := gtime.New(req.CreateTime)

	// Parse first paid time if provided
	var firstPaidTime *gtime.Time
	if len(req.FirstPaidTime) > 0 {
		firstPaidTime = gtime.New(req.FirstPaidTime)
	} else {
		firstPaidTime = gtime.Now()
	}

	// Validate time relationships
	now := gtime.Now()
	utility.Assert(currentPeriodStart.Timestamp() <= now.Timestamp(), "CurrentPeriodStart should earlier then now")
	utility.Assert(currentPeriodEnd.Timestamp() > now.Timestamp(), "CurrentPeriodEnd should later then now")
	utility.Assert(billingCycleAnchor.Timestamp() <= now.Timestamp(), "BillingCycleAnchor should earlier then now")
	if firstPaidTime != nil {
		utility.Assert(firstPaidTime.Timestamp() <= now.Timestamp(), "FirstPaidTime should earlier then now")
	}
	utility.Assert(createTime.Timestamp() <= now.Timestamp(), "CreateTime should earlier then now")
	utility.Assert(currentPeriodStart.Timestamp() >= createTime.Timestamp() && currentPeriodStart.Timestamp() >= billingCycleAnchor.Timestamp(), "currentPeriodStart should later then createTime and billingCycleAnchor")
	utility.Assert(currentPeriodEnd.Timestamp() > currentPeriodStart.Timestamp() && currentPeriodEnd.Timestamp() > billingCycleAnchor.Timestamp(), "currentPeriodEnd should later then currentPeriodStart and billingCycleAnchor")
	if firstPaidTime != nil {
		utility.Assert(currentPeriodEnd.Timestamp() > firstPaidTime.Timestamp(), "currentPeriodEnd should later then firstPaidTime")
	}
	utility.Assert(currentPeriodEnd.Timestamp() > createTime.Timestamp(), "currentPeriodEnd should later then createTime")

	if firstPaidTime == nil {
		firstPaidTime = currentPeriodStart
	}
	// Check if subscription already exists
	existingSubscription := query.GetSubscriptionByExternalSubscriptionId(ctx, req.ExternalSubscriptionId)
	tag := ""

	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}

	if _interface.Context().Get(ctx).MerchantMember != nil {
		tag = fmt.Sprintf("ImportByMember:%d", _interface.Context().Get(ctx).MerchantMember.Id)
		req.Metadata["ImportFrom"] = tag
	} else {
		tag = fmt.Sprintf("ImportByOpenAPI")
		req.Metadata["ImportFrom"] = tag
	}

	if existingSubscription != nil {
		// Check permission to override
		utility.Assert(existingSubscription.Data == tag, "no permission to override, "+existingSubscription.Data)
		utility.Assert(existingSubscription.UserId == user.Id, "no permission to override, user not match")
		var dunningTime = period.GetDunningTimeFromEnd(ctx, utility.MaxInt64(currentPeriodEnd.Timestamp(), 0), plan.Id)
		// Update existing subscription
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().Type:                        consts.SubTypeUniBeeControl,
			dao.Subscription.Columns().Status:                      consts.SubStatusActive,
			dao.Subscription.Columns().Amount:                      totalAmount,
			dao.Subscription.Columns().Currency:                    plan.Currency,
			dao.Subscription.Columns().PlanId:                      plan.Id,
			dao.Subscription.Columns().Quantity:                    req.Quantity,
			dao.Subscription.Columns().GatewayId:                   gatewayId,
			dao.Subscription.Columns().GatewayItemData:             req.Features,
			dao.Subscription.Columns().GatewayDefaultPaymentMethod: "",
			dao.Subscription.Columns().TaxPercentage:               taxPercentage,
			dao.Subscription.Columns().BillingCycleAnchor:          billingCycleAnchor.Timestamp(),
			dao.Subscription.Columns().CurrentPeriodStart:          currentPeriodStart.Timestamp(),
			dao.Subscription.Columns().CurrentPeriodEnd:            currentPeriodEnd.Timestamp(),
			dao.Subscription.Columns().CurrentPeriodStartTime:      currentPeriodStart,
			dao.Subscription.Columns().CurrentPeriodEndTime:        currentPeriodEnd,
			dao.Subscription.Columns().DunningTime:                 dunningTime,
			dao.Subscription.Columns().FirstPaidTime:               firstPaidTime.Timestamp(),
			dao.Subscription.Columns().CreateTime:                  createTime.Timestamp(),
			dao.Subscription.Columns().MetaData:                    utility.MarshalToJsonString(req.Metadata),
		}).Where(dao.Subscription.Columns().Id, existingSubscription.Id).OmitNil().Update()
		utility.AssertError(err, "Update subscription error")

		{
			if len(existingSubscription.LatestInvoiceId) == 0 {
				currentInvoice := &bean.Invoice{
					InvoiceName:                    "SubscriptionCreate",
					BizType:                        consts.BizTypeSubscription,
					ProductName:                    plan.PlanName,
					OriginAmount:                   0,
					TotalAmount:                    0,
					TotalAmountExcludingTax:        0,
					DiscountCode:                   "",
					DiscountAmount:                 0,
					Currency:                       existingSubscription.Currency,
					TaxAmount:                      0,
					SubscriptionAmount:             0,
					SubscriptionAmountExcludingTax: 0,
					Lines: []*bean.InvoiceItemSimplify{{
						Currency:               existingSubscription.Currency,
						OriginAmount:           0,
						Amount:                 0,
						DiscountAmount:         0,
						Tax:                    0,
						AmountExcludingTax:     0,
						TaxPercentage:          0,
						UnitAmountExcludingTax: 0,
						Name:                   plan.PlanName,
						Description:            plan.Description,
						Proration:              false,
						Quantity:               req.Quantity,
						PeriodEnd:              currentPeriodEnd.Timestamp(),
						PeriodStart:            currentPeriodStart.Timestamp(),
						Plan:                   bean.SimplifyPlan(plan),
					}},
					ProrationDate: time.Now().Unix(),
					PeriodStart:   existingSubscription.CurrentPeriodStart,
					PeriodEnd:     existingSubscription.CurrentPeriodEnd,
					Metadata:      req.Metadata,
					CountryCode:   countryCode,
					VatNumber:     vatNumber,
					TaxPercentage: 0,
				}
				invoice, err := service3.CreateProcessingInvoiceForSub(ctx, &service3.CreateProcessingInvoiceForSubReq{
					PlanId:             plan.Id,
					Simplify:           currentInvoice,
					Sub:                existingSubscription,
					GatewayId:          existingSubscription.GatewayId,
					IsSubLatestInvoice: true,
					TimeNow:            currentInvoice.ProrationDate,
				})
				utility.AssertError(err, "Create Latest Invoice Error")
				invoice, err = handler2.MarkInvoiceAsPaidForZeroPayment(ctx, invoice.InvoiceId)
				utility.AssertError(err, "Create Latest Invoice Error")
				timeline.SubscriptionNewTimeline(ctx, invoice)
				sub_update.UpdateUserDefaultSubscriptionForPaymentSuccess(ctx, existingSubscription.UserId, existingSubscription.SubscriptionId)
			}
		}

		subscription3.SendMerchantSubscriptionWebhookBackground(existingSubscription, -10000, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_IMPORT_OVERRIDE, map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()})

		// Log operation
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     existingSubscription.MerchantId,
			Target:         fmt.Sprintf("Subscription(%s)", existingSubscription.SubscriptionId),
			Content:        "ImportOverride",
			UserId:         existingSubscription.UserId,
			SubscriptionId: existingSubscription.SubscriptionId,
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)

		if req.CountryCode != "" && req.CountryCode != user.CountryCode {
			sub_update.UpdateUserCountryCode(ctx, user.Id, countryCode)
		}
		if fmt.Sprintf("%d", gatewayId) != user.GatewayId {
			sub_update.UpdateUserDefaultGatewayForCheckout(ctx, user.Id, gatewayId, "")
		}

		// Get updated subscription detail
		subscriptionDetail, err := detailService.SubscriptionDetail(ctx, existingSubscription.SubscriptionId)
		utility.AssertError(err, "Get subscription detail error")

		if subscriptionDetail != nil {
			sub_update.UpdateUserDefaultSubscriptionForUpdate(ctx, subscriptionDetail.User.Id, subscriptionDetail.Subscription.SubscriptionId)
		}

		return &subscription.ActiveSubscriptionImportRes{
			Subscription: subscriptionDetail,
		}, nil
	} else {
		// Create new subscription
		var dunningTime = period.GetDunningTimeFromEnd(ctx, utility.MaxInt64(currentPeriodEnd.Timestamp(), 0), plan.Id)
		newSubscription := &entity.Subscription{
			SubscriptionId:              utility.CreateSubscriptionId(),
			ExternalSubscriptionId:      req.ExternalSubscriptionId,
			Type:                        consts.SubTypeUniBeeControl,
			UserId:                      user.Id,
			Amount:                      totalAmount,
			Currency:                    plan.Currency,
			MerchantId:                  merchantId,
			PlanId:                      plan.Id,
			Quantity:                    req.Quantity,
			GatewayId:                   gatewayId,
			Status:                      consts.SubStatusActive,
			CurrentPeriodStart:          currentPeriodStart.Timestamp(),
			CurrentPeriodEnd:            currentPeriodEnd.Timestamp(),
			CurrentPeriodStartTime:      currentPeriodStart,
			CurrentPeriodEndTime:        currentPeriodEnd,
			DunningTime:                 dunningTime,
			BillingCycleAnchor:          billingCycleAnchor.Timestamp(),
			FirstPaidTime:               firstPaidTime.Timestamp(),
			CreateTime:                  createTime.Timestamp(),
			CountryCode:                 user.CountryCode,
			VatNumber:                   user.VATNumber,
			TaxPercentage:               taxPercentage,
			GatewaySubscriptionId:       req.ExternalSubscriptionId,
			GatewayItemData:             req.Features,
			Data:                        tag,
			CurrentPeriodPaid:           1,
			GatewayDefaultPaymentMethod: "",
			MetaData:                    utility.MarshalToJsonString(req.Metadata),
		}

		result, err := dao.Subscription.Ctx(ctx).Data(newSubscription).OmitNil().Insert(newSubscription)
		utility.AssertError(err, "Save subscription error")
		id, err := result.LastInsertId()
		utility.AssertError(err, "Get last insert id error")
		newSubscription.Id = uint64(id)

		{
			currentInvoice := &bean.Invoice{
				InvoiceName:                    "SubscriptionCreate",
				BizType:                        consts.BizTypeSubscription,
				ProductName:                    plan.PlanName,
				OriginAmount:                   0,
				TotalAmount:                    0,
				TotalAmountExcludingTax:        0,
				DiscountCode:                   "",
				DiscountAmount:                 0,
				Currency:                       newSubscription.Currency,
				TaxAmount:                      0,
				SubscriptionAmount:             0,
				SubscriptionAmountExcludingTax: 0,
				Lines: []*bean.InvoiceItemSimplify{{
					Currency:               newSubscription.Currency,
					OriginAmount:           0,
					Amount:                 0,
					DiscountAmount:         0,
					Tax:                    0,
					AmountExcludingTax:     0,
					TaxPercentage:          0,
					UnitAmountExcludingTax: 0,
					Name:                   plan.PlanName,
					Description:            plan.Description,
					Proration:              false,
					Quantity:               req.Quantity,
					PeriodEnd:              currentPeriodEnd.Timestamp(),
					PeriodStart:            currentPeriodStart.Timestamp(),
					Plan:                   bean.SimplifyPlan(plan),
				}},
				ProrationDate: time.Now().Unix(),
				PeriodStart:   newSubscription.CurrentPeriodStart,
				PeriodEnd:     newSubscription.CurrentPeriodEnd,
				Metadata:      req.Metadata,
				CountryCode:   countryCode,
				VatNumber:     vatNumber,
				TaxPercentage: 0,
			}
			invoice, err := service3.CreateProcessingInvoiceForSub(ctx, &service3.CreateProcessingInvoiceForSubReq{
				PlanId:             plan.Id,
				Simplify:           currentInvoice,
				Sub:                newSubscription,
				GatewayId:          newSubscription.GatewayId,
				IsSubLatestInvoice: true,
				TimeNow:            currentInvoice.ProrationDate,
			})
			utility.AssertError(err, "Create Latest Invoice Error")
			invoice, err = handler2.MarkInvoiceAsPaidForZeroPayment(ctx, invoice.InvoiceId)
			utility.AssertError(err, "Create Latest Invoice Error")
			timeline.SubscriptionNewTimeline(ctx, invoice)
			sub_update.UpdateUserDefaultSubscriptionForPaymentSuccess(ctx, newSubscription.UserId, newSubscription.SubscriptionId)
		}

		subscription3.SendMerchantSubscriptionWebhookBackground(newSubscription, -10000, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_IMPORT_CREATED, map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()})

		// Log operation
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     newSubscription.MerchantId,
			Target:         fmt.Sprintf("Subscription(%s)", newSubscription.SubscriptionId),
			Content:        "ImportNew",
			UserId:         newSubscription.UserId,
			SubscriptionId: newSubscription.SubscriptionId,
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)

		if req.CountryCode != "" && req.CountryCode != user.CountryCode {
			sub_update.UpdateUserCountryCode(ctx, user.Id, countryCode)
		}
		if fmt.Sprintf("%d", gatewayId) != user.GatewayId {
			sub_update.UpdateUserDefaultGatewayForCheckout(ctx, user.Id, gatewayId, "")
		}

		// Get subscription detail
		subscriptionDetail, err := detailService.SubscriptionDetail(ctx, newSubscription.SubscriptionId)
		utility.AssertError(err, "Get subscription detail error")

		if subscriptionDetail != nil {
			sub_update.UpdateUserDefaultSubscriptionForUpdate(ctx, subscriptionDetail.User.Id, subscriptionDetail.Subscription.SubscriptionId)
		}

		return &subscription.ActiveSubscriptionImportRes{
			Subscription: subscriptionDetail,
		}, nil
	}
}
