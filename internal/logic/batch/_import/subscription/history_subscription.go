package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"strconv"
	"strings"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/plan/period"
	user2 "unibee/internal/logic/user"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

type TaskHistorySubscriptionImport struct {
}

func (t TaskHistorySubscriptionImport) TaskName() string {
	return "HistorySubscriptionImport"
}

func (t TaskHistorySubscriptionImport) TemplateVersion() string {
	return "v1"
}

func (t TaskHistorySubscriptionImport) TemplateHeader() interface{} {
	return &ImportHistorySubscriptionEntity{
		ExternalSubscriptionId: "exampleSubscriptionId",
		ExternalUserId:         "exampleUserId",
		ExternalPlanId:         "examplePlanId",
		Amount:                 "10.00",
		Quantity:               "1",
		CountryCode:            "EE",
		TaxPercentage:          "1000",
		Gateway:                "stripe",
		CurrentPeriodStart:     "2024-05-13 06:19:27",
		CurrentPeriodEnd:       "2024-06-13 06:19:27",
	}
}

func (t TaskHistorySubscriptionImport) ImportRow(ctx context.Context, task *entity.MerchantBatchTask, row map[string]string) (interface{}, error) {
	var err error
	target := &ImportHistorySubscriptionEntity{
		ExternalSubscriptionId: fmt.Sprintf("%s", row["ExternalSubscriptionId"]),
		ExternalUserId:         fmt.Sprintf("%s", row["ExternalUserId"]),
		ExternalPlanId:         fmt.Sprintf("%s", row["ExternalPlanId"]),
		Amount:                 fmt.Sprintf("%s", row["Amount"]),
		Quantity:               fmt.Sprintf("%s", row["Quantity"]),
		CountryCode:            fmt.Sprintf("%s", row["CountryCode"]),
		TaxPercentage:          fmt.Sprintf("%s", row["TaxPercentage"]),
		Gateway:                fmt.Sprintf("%s", row["Gateway"]),
		CurrentPeriodStart:     fmt.Sprintf("%s", row["CurrentPeriodStart"]),
		CurrentPeriodEnd:       fmt.Sprintf("%s", row["CurrentPeriodEnd"]),
	}
	tag := fmt.Sprintf("ImportBy%v", task.MemberId)
	if len(target.ExternalSubscriptionId) == 0 {
		return target, gerror.New("Error, ExternalSubscriptionId is blank")
	}
	target.CountryCode = strings.ToUpper(target.CountryCode)
	if len(target.CountryCode) > 0 {
		err = utility.ValidateCountryCode(target.CountryCode)
		if err != nil {
			return target, gerror.New(fmt.Sprintf("Error, CountryCode is invalid, %s", err.Error()))
		}
	}
	// data prepare
	user, err := user2.QueryOrCreateUser(ctx, &user2.NewUserInternalReq{
		ExternalUserId: target.ExternalUserId,
		Email:          target.Email,
		CountryCode:    target.CountryCode,
		MerchantId:     _interface.GetMerchantId(ctx),
	})
	if err != nil {
		return target, gerror.Newf("QueryOrCreateUser,error:%s", err.Error())
	}
	if user == nil {
		return target, gerror.New("Error, can't find user by ExternalUserId")
	}

	// Get tax percentage from request instead of user
	var taxPercentage int64 = 0
	if len(target.TaxPercentage) > 0 {
		taxPercentageInt, err := strconv.ParseInt(target.TaxPercentage, 10, 64)
		if err != nil {
			return target, gerror.Newf("Invalid TaxPercentage,error:%s", err.Error())
		}
		taxPercentage = taxPercentageInt
	}

	if len(target.ExternalPlanId) == 0 {
		return target, gerror.New("Error, ExternalPlanId is blank")
	}
	plan := query.GetPlanByExternalPlanId(ctx, task.MerchantId, target.ExternalPlanId)
	if plan == nil {
		return target, gerror.New("Error, can't find plan by ExternalPlanId")
	}
	utility.Assert(plan.Status != consts.PlanStatusEditable, "plan status should not editable")
	if len(target.Amount) == 0 {
		return target, gerror.New("Error, Amount is blank")
	}
	amountFloat, err := strconv.ParseFloat(target.Amount, 64)
	if err != nil {
		return target, gerror.Newf("Invalid Amount,error:%s", err.Error())
	}
	amount := int64(amountFloat * 100)
	if amount <= 0 {
		return target, gerror.New("Invalid Amount, should greater than 0")
	}
	if len(target.Gateway) == 0 {
		return target, gerror.New("Error, Gateway is blank")
	}
	var gatewayId uint64 = 0
	gatewayImpl := api.GatewayNameMapping[target.Gateway]
	if gatewayImpl == nil {
		return target, gerror.New("Error, Invalid Gateway, should be one of " + strings.Join(api.ExportGatewaySetupListKeys(), "|"))
	}
	gateway := query.GetDefaultGatewayByGatewayName(ctx, task.MerchantId, target.Gateway)
	if gateway == nil {
		return target, gerror.New("Error, gateway need setup")
	}
	gatewayId = gateway.Id
	quantity, _ := strconv.ParseInt(target.Quantity, 10, 64)
	if quantity == 0 {
		quantity = 1
	}
	if len(target.CurrentPeriodStart) == 0 {
		return target, gerror.New("Error, CurrentPeriodStart is blank")
	}
	currentPeriodStart := gtime.New(target.CurrentPeriodStart)
	if len(target.CurrentPeriodEnd) == 0 {
		return target, gerror.New("Error, CurrentPeriodEnd is blank")
	}
	currentPeriodEnd := gtime.New(target.CurrentPeriodEnd)

	// data verify
	{
		if currentPeriodStart.Timestamp() >= gtime.Now().Timestamp() {
			return target, gerror.New("Error, CurrentPeriodStart should earlier then now")
		}
		if currentPeriodEnd.Timestamp() >= gtime.Now().Timestamp() {
			return target, gerror.New("Error, CurrentPeriodEnd should earlier then now")
		}
		if currentPeriodEnd.Timestamp() <= currentPeriodStart.Timestamp() {
			return target, gerror.New("Error,currentPeriodEnd should later then currentPeriodStart")
		}
	}
	var dunningTime = period.GetDunningTimeFromEnd(ctx, utility.MaxInt64(currentPeriodEnd.Timestamp(), 0), plan.Id)
	one := query.GetSubscriptionByExternalSubscriptionId(ctx, target.ExternalSubscriptionId)
	if one != nil {
		if one.UserId != user.Id {
			return target, gerror.New("Error, user not match")
		}
		if one.Data != tag {
			return target, gerror.New("Error, no permission to override," + one.Data)
		}
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().Type:                   consts.SubTypeUniBeeControl,
			dao.Subscription.Columns().Amount:                 amount,
			dao.Subscription.Columns().Currency:               plan.Currency,
			dao.Subscription.Columns().PlanId:                 plan.Id,
			dao.Subscription.Columns().Quantity:               quantity,
			dao.Subscription.Columns().GatewayId:              gatewayId,
			dao.Subscription.Columns().TaxPercentage:          taxPercentage,
			dao.Subscription.Columns().CountryCode:            target.CountryCode,
			dao.Subscription.Columns().BillingCycleAnchor:     currentPeriodStart.Timestamp(),
			dao.Subscription.Columns().CurrentPeriodStart:     currentPeriodStart.Timestamp(),
			dao.Subscription.Columns().CurrentPeriodEnd:       currentPeriodEnd.Timestamp(),
			dao.Subscription.Columns().CurrentPeriodStartTime: currentPeriodStart,
			dao.Subscription.Columns().CurrentPeriodEndTime:   currentPeriodEnd,
			dao.Subscription.Columns().DunningTime:            dunningTime,
			dao.Subscription.Columns().FirstPaidTime:          currentPeriodStart.Timestamp(),
			dao.Subscription.Columns().CreateTime:             currentPeriodStart.Timestamp(),
		}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return target, gerror.Newf("Override subscription,error:%s", err.Error())
		}
	} else {
		one = &entity.Subscription{
			Type:                   consts.SubTypeUniBeeControl,
			SubscriptionId:         utility.CreateSubscriptionId(),
			ExternalSubscriptionId: target.ExternalSubscriptionId,
			UserId:                 user.Id,
			Amount:                 amount,
			Currency:               plan.Currency,
			MerchantId:             task.MerchantId,
			PlanId:                 plan.Id,
			Quantity:               quantity,
			GatewayId:              gatewayId,
			Status:                 consts.SubStatusExpired,
			CurrentPeriodStart:     currentPeriodStart.Timestamp(),
			CurrentPeriodEnd:       currentPeriodEnd.Timestamp(),
			DunningTime:            dunningTime,
			CurrentPeriodStartTime: currentPeriodStart,
			CurrentPeriodEndTime:   currentPeriodEnd,
			BillingCycleAnchor:     currentPeriodStart.Timestamp(),
			FirstPaidTime:          currentPeriodStart.Timestamp(),
			CreateTime:             currentPeriodStart.Timestamp(),
			CountryCode:            target.CountryCode,
			VatNumber:              user.VATNumber,
			TaxPercentage:          taxPercentage,
			GatewaySubscriptionId:  target.ExternalSubscriptionId,
			Data:                   tag,
			CurrentPeriodPaid:      1,
		}
		_, err = dao.Subscription.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			return target, gerror.Newf("Save subscription,error:%s", err.Error())
		}
	}
	uniqueId := fmt.Sprintf("%s-%v-%v-%v", one.ExternalSubscriptionId, one.PlanId, one.CurrentPeriodStart, one.CurrentPeriodEnd)
	timeline := query.GetSubscriptionTimeLineByUniqueId(ctx, uniqueId)
	if timeline != nil {
		return target, gerror.New("same history record exist:" + uniqueId)
	}
	// to start import timeline
	timeline = &entity.SubscriptionTimeline{
		MerchantId:      one.MerchantId,
		UserId:          one.UserId,
		SubscriptionId:  one.SubscriptionId,
		UniqueId:        uniqueId,
		Currency:        one.Currency,
		PlanId:          one.PlanId,
		Quantity:        one.Quantity,
		AddonData:       one.AddonData,
		Status:          consts.SubTimeLineStatusFinished,
		GatewayId:       one.GatewayId,
		PeriodStart:     one.CurrentPeriodStart,
		PeriodEnd:       one.CurrentPeriodEnd,
		PeriodStartTime: gtime.NewFromTimeStamp(one.CurrentPeriodStart),
		PeriodEndTime:   gtime.NewFromTimeStamp(one.CurrentPeriodEnd),
		CreateTime:      gtime.Now().Timestamp(),
	}

	result, err := dao.SubscriptionTimeline.Ctx(ctx).Data(timeline).OmitNil().Insert(timeline)
	if err != nil {
		return target, gerror.Newf("Save history,error:%s", err.Error())
	}
	id, err := result.LastInsertId()
	utility.AssertError(err, "Save history error")
	timeline.Id = uint64(id)

	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("SubscriptionHistory(%v)", timeline.Id),
		Content:        "ImportNew",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)

	return target, err
}

type ImportHistorySubscriptionEntity struct {
	ExternalSubscriptionId string `json:"ExternalSubscriptionId"    comment:"Required, The external id of subscription"     `
	ExternalPlanId         string `json:"ExternalPlanId"   comment:"Required, The external id of plan, plan should created at first"   `
	Email                  string `json:"Email"  comment:"The email of user, one of Email or ExternalUserId is required" `
	ExternalUserId         string `json:"ExternalUserId"    comment:"The external id of user, one of Email or ExternalUserId is required "    `
	Amount                 string `json:"Amount"        comment:"Required, the recurring amount of subscription, em. 19.99 = 19.99 USD"     `
	Quantity               string `json:"Quantity"      comment:"the quantity of plan, default 1 if not provided "        `
	CountryCode            string `json:"CountryCode"    comment:"Required, The country code of subscription, Tax is applied based on this country code."  `
	TaxPercentage          string `json:"TaxPercentage" comment:"The TaxPercentage of subscription, 1000 = 10%"`
	Gateway                string `json:"Gateway" comment:"Required, should one of stripe|paypal|wire_transfer|changelly "           `
	CurrentPeriodStart     string `json:"CurrentPeriodStart" comment:"Required, UTC time, the current period start time of subscription, format '2006-01-02 15:04:05'"`
	CurrentPeriodEnd       string `json:"CurrentPeriodEnd"   comment:"Required, UTC time, the current period end time of subscription, format '2006-01-02 15:04:05'"`
}
