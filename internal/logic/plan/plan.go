package plan

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"unibee/api/bean"
	"unibee/api/merchant/plan"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/currency"
	"unibee/internal/logic/metric"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"
)

var PlanIntervals = []string{"day", "month", "year", "week"}
var PlanIntervalLength = map[string]int{
	"day":   1,
	"week":  7,
	"month": 30,
	"year":  365,
}

func PlanPublish(ctx context.Context, planId uint64) (err error) {
	utility.Assert(planId > 0, "invalid planId")
	one := query.GetPlanById(ctx, planId)
	utility.Assert(one.Status == consts.PlanStatusActive, "plan not activate")
	PlanOrAddonIntervalVerify(ctx, planId)
	_, err = dao.Plan.Ctx(ctx).Data(g.Map{
		dao.Plan.Columns().PublishStatus: consts.PlanPublishStatusPublished,
		dao.Plan.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.Plan.Columns().Id, planId).Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Plan(%v)", one.Id),
		Content:        "Publish",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         one.Id,
		DiscountCode:   "",
	}, err)
	if err != nil {
		return err
	}
	return nil
}

func PlanUnPublish(ctx context.Context, planId uint64) (err error) {
	utility.Assert(planId > 0, "invalid planId")
	one := query.GetPlanById(ctx, planId)
	utility.Assert(one.Status == consts.PlanStatusActive, "plan not activate")
	_, err = dao.Plan.Ctx(ctx).Data(g.Map{
		dao.Plan.Columns().PublishStatus: consts.PlanPublishStatusUnPublished,
		dao.Plan.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.Plan.Columns().Id, planId).Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Plan(%v)", one.Id),
		Content:        "UnPublish",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         one.Id,
		DiscountCode:   "",
	}, err)
	if err != nil {
		return err
	}
	return nil
}

type PlanInternalReq struct {
	ExternalPlanId        string                               `json:"externalPlanId" dc:"ExternalPlanId"`
	MerchantId            uint64                               `json:"merchantId" dc:"MerchantId" `
	PlanId                uint64                               `json:"planId" dc:"PlanId" `
	PlanName              string                               `json:"planName" dc:"Plan Name"    `
	InternalName          string                               `json:"internalName"              description:""` //
	Amount                int64                                `json:"amount"   dc:"Plan CaptureAmount"  `
	Currency              string                               `json:"currency"   dc:"Plan Currency"`
	IntervalUnit          string                               `json:"intervalUnit" dc:"Plan Interval Unit，em: day|month|year|week"`
	IntervalCount         int                                  `json:"intervalCount"  dc:"Number Of IntervalUnit，em: day|month|year|week"`
	Description           string                               `json:"description"  dc:"Description"`
	Type                  int                                  `json:"type"  d:"1"  dc:"Default 1，,1-main plan，2-addon plan, 3-onetime plan" `
	ProductName           string                               `json:"productName" dc:"Default Copy PlanName"  `
	ProductDescription    string                               `json:"productDescription" dc:"Default Copy Description" `
	ImageUrl              string                               `json:"imageUrl"    dc:"ImageUrl,Start With: http" `
	HomeUrl               string                               `json:"homeUrl"    dc:"HomeUrl,Start With: http"  `
	AddonIds              []int64                              `json:"addonIds"  dc:"Plan Ids Of Recurring Addon Type" `
	OnetimeAddonIds       []int64                              `json:"onetimeAddonIds"  dc:"Plan Ids Of Onetime Addon Type" `
	MetricLimits          []*bean.PlanMetricLimitParam         `json:"metricLimits"  dc:"Plan's MetricLimit List" `
	MetricMeteredCharge   []*bean.PlanMetricMeteredChargeParam `json:"metricMeteredCharge"  dc:"Plan's MetricMeteredCharge" `
	MetricRecurringCharge []*bean.PlanMetricMeteredChargeParam `json:"metricRecurringCharge"  dc:"Plan's MetricRecurringCharge" `
	GasPayer              string                               `json:"gasPayer" dc:"who pay the gas for crypto payment, merchant|user"`
	Metadata              map[string]interface{}               `json:"metadata" dc:"Metadata，Map"`
	TrialAmount           int64                                `json:"trialAmount"                description:"price of trial period"` // price of trial period
	TrialDurationTime     int64                                `json:"trialDurationTime"         description:"duration of trial"`      // duration of trial
	TrialDemand           string                               `json:"trialDemand"               description:"demand of trial, example, paymentMethod, payment method will ask for subscription trial start"`
	CancelAtTrialEnd      int                                  `json:"cancelAtTrialEnd"          description:"whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription"` // whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription
	ProductId             int64                                `json:"productId"   dc:"Id of product which plan to linked" `
}

func MetricPlanChargeValidation(metricPlanCharges []*bean.PlanMetricMeteredChargeParam) error {
	for _, metricPlanCharge := range metricPlanCharges {
		if metricPlanCharge.MetricId <= 0 {
			return gerror.New("metric id should not less than 0")
		}
		if metricPlanCharge.ChargeType != 0 && metricPlanCharge.ChargeType != 1 {
			return gerror.New("charge type should be one of 0-standard pricing, 1-graduated pricing")
		}
		if metricPlanCharge.ChargeType == 0 {
			if metricPlanCharge.StandardAmount < 0 {
				return gerror.New("standard amount should not be negative")
			}
			if metricPlanCharge.StandardStartValue < 0 {
				return gerror.New("standard start value should not be negative")
			}
		} else if metricPlanCharge.ChargeType == 1 {
			var lastEnd int64 = 0
			for _, step := range metricPlanCharge.GraduatedAmounts {
				if step.EndValue > 0 && step.EndValue <= lastEnd {
					return gerror.New("end value should be greater than last end value")
				}
				if step.PerAmount < 0 {
					return gerror.New("per amount should not be negative")
				}
				if step.FlatAmount < 0 {
					return gerror.New("flat amount should not be negative")
				}
				lastEnd = step.EndValue
			}
			if lastEnd > 0 {
				return gerror.New("The last EndValue should the infinity value")
			}
		}
	}
	return nil
}

func PlanCreate(ctx context.Context, req *PlanInternalReq) (one *entity.Plan, err error) {
	utility.Assert(req.MerchantId > 0, "merchantId invalid")
	intervals := []string{"day", "month", "year", "week"}
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.Amount >= 0, "amount value should >= 0")
	utility.Assert(len(req.PlanName) > 0, "plan name should not blank")
	utility.Assert(currency.IsFiatCurrencySupport(req.Currency), "currency not support")
	// product check and update
	if req.ProductId > 0 {
		product := query.GetProductById(ctx, uint64(req.ProductId), req.MerchantId)
		utility.Assert(product != nil && product.Status == 1, "product not found or product not activate")
	}

	if len(req.ExternalPlanId) > 0 {
		utility.Assert(query.GetPlanByExternalPlanId(ctx, req.MerchantId, req.ExternalPlanId) == nil, "Same ExternalPlanId Exist")
	}

	if len(req.GasPayer) > 0 {
		utility.Assert(strings.Contains("merchant|user", req.GasPayer), "gasPayer should one of merchant|user")
	}

	if req.Type != consts.PlanTypeMain {
		if req.Type == consts.PlanTypeRecurringAddon {
			utility.Assert(req.TrialDurationTime == 0, "Trial not available for addon")
			utility.Assert(req.TrialAmount == 0, "Trial not available for addon")
			utility.Assert(req.TrialDemand == "", "Trial not available for addon")
		}
		utility.Assert(len(req.MetricMeteredCharge) == 0, "Metric metered charge not available for addon")
		utility.Assert(len(req.MetricRecurringCharge) == 0, "Metric recurring charge not available for addon")
	}

	//check metricLimitList
	if len(req.MetricLimits) > 0 {
		for _, ml := range req.MetricLimits {
			utility.Assert(ml.MetricId > 0, "invalid metricId")
			utility.Assert(ml.MetricLimit > 0, "invalid MetricLimit")
			me := query.GetMerchantMetric(ctx, ml.MetricId)
			utility.Assert(me != nil, "metric not found")
			utility.Assert(me.Type == metric.MetricTypeLimitMetered, "metric type invalid")
		}
	}
	merchantInfo := query.GetMerchantById(ctx, req.MerchantId)
	if len(req.ImageUrl) == 0 {
		req.ImageUrl = merchantInfo.CompanyLogo
	}
	if len(req.HomeUrl) == 0 {
		req.HomeUrl = merchantInfo.HomeUrl
	}
	utility.Assert(merchantInfo != nil, "merchant not found")
	utility.Assert(req.Type == consts.PlanTypeMain || req.Type == consts.PlanTypeRecurringAddon || req.Type == consts.PlanTypeOnetime, "type should be 1|2｜3")
	if req.Type != consts.PlanTypeOnetime {
		utility.Assert(utility.StringContainsElement(intervals, strings.ToLower(req.IntervalUnit)), "IntervalUnit Error， must one of day｜month｜year｜week\"")
		utility.Assert(req.IntervalCount > 0, "IntervalCount should > 0")
		if strings.ToLower(req.IntervalUnit) == "day" {
			utility.Assert(req.IntervalCount <= 365, "IntervalCount Must Lower Then 365 While IntervalUnit is day")
		} else if strings.ToLower(req.IntervalUnit) == "month" {
			utility.Assert(req.IntervalCount <= 12, "IntervalCount Must Lower Then 12 While IntervalUnit is month")
		} else if strings.ToLower(req.IntervalUnit) == "year" {
			utility.Assert(req.IntervalCount <= 1, "IntervalCount Must Lower Then 2 While IntervalUnit is year")
		} else if strings.ToLower(req.IntervalUnit) == "week" {
			utility.Assert(req.IntervalCount <= 52, "IntervalCount Must Lower Then 52 While IntervalUnit is week")
		}
	}
	if req.IntervalCount < 1 {
		req.IntervalCount = 1
	}

	if len(req.ProductName) == 0 {
		req.ProductName = req.PlanName
	}
	if len(req.ProductDescription) == 0 {
		req.ProductDescription = req.Description
	}

	if len(req.AddonIds) > 0 {
		var allAddonList []*entity.Plan
		err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, req.AddonIds).OmitEmpty().Scan(&allAddonList)
		for _, addonPlan := range allAddonList {
			utility.Assert(addonPlan.Type == consts.PlanTypeRecurringAddon, fmt.Sprintf("plan not recurring addon type, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("addon not published status, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Currency == req.Currency, fmt.Sprintf("addon currency not match plan's currency, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.ProductId == req.ProductId, fmt.Sprintf("addon product not match plan's product, id:%d", addonPlan.Id))
		}
	}

	if len(req.OnetimeAddonIds) > 0 {
		var allAddonList []*entity.Plan
		err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, req.OnetimeAddonIds).OmitEmpty().Scan(&allAddonList)
		for _, addonPlan := range allAddonList {
			utility.Assert(addonPlan.Type == consts.PlanTypeOnetime, fmt.Sprintf("plan not onetime addon type, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("add plan not published status, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Currency == req.Currency, fmt.Sprintf("add plan currency not match plan's currency, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.ProductId == req.ProductId, fmt.Sprintf("addon product not match plan's product, id:%d", addonPlan.Id))
		}
	}

	if len(req.MetricMeteredCharge) > 0 {
		utility.AssertError(MetricPlanChargeValidation(req.MetricMeteredCharge), "Usage-based Metered charges validation failed")
	}

	if len(req.MetricRecurringCharge) > 0 {
		utility.AssertError(MetricPlanChargeValidation(req.MetricRecurringCharge), "Usage-based Recurring charges validation failed")
	}

	utility.Assert(req.TrialDemand == "" || req.TrialDemand == "paymentMethod", "Demand of trial should be paymentMethod or not")

	one = &entity.Plan{
		ExternalPlanId:         req.ExternalPlanId,
		CompanyId:              merchantInfo.CompanyId,
		MerchantId:             req.MerchantId,
		PlanName:               req.PlanName,
		InternalName:           req.InternalName,
		Amount:                 req.Amount,
		Currency:               strings.ToUpper(req.Currency),
		IntervalUnit:           strings.ToLower(req.IntervalUnit),
		IntervalCount:          req.IntervalCount,
		Type:                   req.Type,
		Description:            req.Description,
		ImageUrl:               req.ImageUrl,
		HomeUrl:                req.HomeUrl,
		BindingAddonIds:        utility.IntListToString(req.AddonIds),
		BindingOnetimeAddonIds: utility.IntListToString(req.OnetimeAddonIds),
		Status:                 consts.PlanStatusEditable,
		CreateTime:             gtime.Now().Timestamp(),
		MetaData:               utility.MarshalToJsonString(req.Metadata),
		GasPayer:               req.GasPayer,
		TrialDurationTime:      req.TrialDurationTime,
		TrialAmount:            req.TrialAmount,
		TrialDemand:            req.TrialDemand,
		CancelAtTrialEnd:       req.CancelAtTrialEnd,
		PublishStatus:          consts.PlanPublishStatusUnPublished,
		ProductId:              req.ProductId,
		MetricCharge: utility.MarshalToJsonString(&bean.MetricPlanBindingEntity{
			MetricLimits:          req.MetricLimits,
			MetricMeteredCharge:   req.MetricMeteredCharge,
			MetricRecurringCharge: req.MetricRecurringCharge,
		}),
	}
	result, err := dao.Plan.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`PlanCreate record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	if len(req.MetricLimits) > 0 {
		err = metric.BulkMetricLimitPlanBindingReplace(ctx, one, req.MetricLimits)
		if err != nil {
			return nil, gerror.Newf(`BulkMetricLimitPlanBindingReplace %s`, err)
		}
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Plan(%v)", one.Id),
		Content:        "New",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         one.Id,
		DiscountCode:   "",
	}, err)
	if req.AddonIds != nil && len(req.AddonIds) > 0 {
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("AddonIds(%s)", utility.IntListToString(req.AddonIds)),
			Content:        "OverrideRecurringAddons",
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         one.Id,
			DiscountCode:   "",
		}, err)
	}
	if req.OnetimeAddonIds != nil && len(req.OnetimeAddonIds) > 0 {
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("OneTimeAddonIds(%s)", utility.IntListToString(req.OnetimeAddonIds)),
			Content:        "OverrideOneTimeAddons",
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         one.Id,
			DiscountCode:   "",
		}, err)
	}
	return one, nil
}

type EditInternalReq struct {
	MerchantId            uint64                                `json:"merchantId" dc:"MerchantId" `
	PlanId                uint64                                `json:"planId" dc:"PlanId" v:"required"`
	ExternalPlanId        *string                               `json:"externalPlanId" dc:"ExternalPlanId"`
	PlanName              *string                               `json:"planName" dc:"Plan Name"   v:"required" `
	InternalName          *string                               `json:"internalName"              description:""` //
	Amount                *int64                                `json:"amount"   dc:"Plan CaptureAmount"   v:"required" `
	Currency              *string                               `json:"currency"   dc:"Plan Currency" v:"required" `
	IntervalUnit          *string                               `json:"intervalUnit" dc:"Plan Interval Unit，em: day|month|year|week"`
	IntervalCount         *int                                  `json:"intervalCount"  dc:"Number Of IntervalUnit" `
	Description           *string                               `json:"description"  dc:"Description"`
	ProductName           *string                               `json:"productName" dc:"Default Copy PlanName"  `
	ProductDescription    *string                               `json:"productDescription" dc:"Default Copy Description" `
	ImageUrl              *string                               `json:"imageUrl"    dc:"ImageUrl,Start With: http" `
	HomeUrl               *string                               `json:"homeUrl"    dc:"HomeUrl,Start With: http"  `
	AddonIds              []int64                               `json:"addonIds"  dc:"Plan Ids Of Recurring Addon Type" `
	OnetimeAddonIds       []int64                               `json:"onetimeAddonIds"  dc:"Plan Ids Of Onetime Addon Type" `
	MetricLimits          []*bean.PlanMetricLimitParam          `json:"metricLimits"  dc:"Plan's MetricLimit List" `
	MetricMeteredCharge   *[]*bean.PlanMetricMeteredChargeParam `json:"metricMeteredCharge"  dc:"Plan's MetricMeteredCharge" `
	MetricRecurringCharge *[]*bean.PlanMetricMeteredChargeParam `json:"metricRecurringCharge"  dc:"Plan's MetricRecurringCharge" `
	GasPayer              *string                               `json:"gasPayer" dc:"who pay the gas for crypto payment, merchant|user"`
	Metadata              *map[string]interface{}               `json:"metadata" dc:"Metadata，Map"`
	TrialAmount           *int64                                `json:"trialAmount"                description:"price of trial period"` // price of trial period
	TrialDurationTime     *int64                                `json:"trialDurationTime"         description:"duration of trial"`      // duration of trial
	TrialDemand           *string                               `json:"trialDemand"               description:"demand of trial, example, paymentMethod, payment method will ask for subscription trial start"`
	CancelAtTrialEnd      *int                                  `json:"cancelAtTrialEnd"          description:"whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription"` // whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription
	ProductId             *int64                                `json:"productId"   dc:"Id of product which plan to linked" `
}

func PlanEdit(ctx context.Context, req *EditInternalReq) (one *entity.Plan, err error) {
	utility.Assert(req != nil, "Req not found")
	utility.Assert(req.PlanId > 0, "PlanId should > 0")
	one = query.GetPlanById(ctx, req.PlanId)
	utility.Assert(one != nil, fmt.Sprintf("plan not found, id:%d", req.PlanId))
	utility.Assert(one.MerchantId == req.MerchantId, "Merchant not match")

	var metricPlanCharge = &bean.MetricPlanBindingEntity{}
	if len(one.MetricCharge) > 0 {
		_ = utility.UnmarshalFromJsonString(one.MetricCharge, &metricPlanCharge)
	}
	// product check and update
	if req.ProductId != nil && *req.ProductId > 0 {
		product := query.GetProductById(ctx, uint64(*req.ProductId), req.MerchantId)
		utility.Assert(product != nil && product.Status == 1, "product not found or product not activate")
	}

	if req.ExternalPlanId != nil && len(*req.ExternalPlanId) > 0 && strings.Compare(one.ExternalPlanId, *req.ExternalPlanId) != 0 {
		utility.Assert(query.GetPlanByExternalPlanId(ctx, req.MerchantId, *req.ExternalPlanId) == nil, "Same ExternalPlanId Exist")
	}

	if one.Status == consts.PlanStatusActive {
		utility.Assert(req.Amount == nil, "Amount is not editable as plan is active")
		utility.Assert(req.Currency == nil, "Currency is not editable as plan is active")
		utility.Assert(req.IntervalUnit == nil, "IntervalUint is not editable as plan is active")
		utility.Assert(req.IntervalCount == nil, "IntervalCount is not editable as plan is active")
	} else {
		if req.Amount != nil {
			utility.Assert(*req.Amount >= 0, "Amount value should >= 0")
		}
		if req.Currency != nil && len(*req.Currency) > 0 {
			utility.Assert(currency.IsFiatCurrencySupport(*req.Currency), "Currency not support")
			one.Currency = *req.Currency
		}

		if req.IntervalCount != nil && *req.IntervalCount < 1 {
			req.IntervalCount = unibee.Int(1)
		}

		if one.Type != consts.PlanTypeOnetime {
			if req.IntervalUnit != nil {
				intervals := []string{"day", "month", "year", "week"}
				utility.Assert(utility.StringContainsElement(intervals, strings.ToLower(*req.IntervalUnit)), "IntervalUnit Error， must one of day｜month｜year｜week\"")
				if req.IntervalCount != nil {
					if strings.ToLower(*req.IntervalUnit) == "day" {
						utility.Assert(*req.IntervalCount <= 365, "IntervalCount Must Lower Then 365 While IntervalUnit is day")
					} else if strings.ToLower(*req.IntervalUnit) == "month" {
						utility.Assert(*req.IntervalCount <= 12, "IntervalCount Must Lower Then 12 While IntervalUnit is month")
					} else if strings.ToLower(*req.IntervalUnit) == "year" {
						utility.Assert(*req.IntervalCount <= 1, "IntervalCount Must Lower Then 2 While IntervalUnit is year")
					} else if strings.ToLower(*req.IntervalUnit) == "week" {
						utility.Assert(*req.IntervalCount <= 52, "IntervalCount Must Lower Then 52 While IntervalUnit is week")
					}
				}
			} else {
				utility.Assert(req.IntervalCount == nil, "IntervalCount can not edit without IntervalUnit")
			}
		}

	}

	if one.Type != consts.PlanTypeMain {
		if one.Type == consts.PlanTypeRecurringAddon {
			utility.Assert(req.TrialDurationTime == nil, "Trial not available for addon")
			utility.Assert(req.TrialAmount == nil, "Trial not available for addon")
			utility.Assert(req.TrialDemand == nil, "Trial not available for addon")
		}
		utility.Assert(req.MetricMeteredCharge == nil || len(*req.MetricMeteredCharge) == 0, "Metric metered charge not available for addon")
		utility.Assert(req.MetricRecurringCharge == nil || len(*req.MetricRecurringCharge) == 0, "Metric recurring charge not available for addon")
	}

	if req.PlanName != nil {
		utility.Assert(len(*req.PlanName) > 0, "Plan name should not blank")
	}

	if req.GasPayer != nil && len(*req.GasPayer) > 0 {
		utility.Assert(strings.Contains("merchant|user", *req.GasPayer), "GasPayer should one of merchant|user")
	}

	//check metricLimitList
	if len(req.MetricLimits) > 0 {
		for _, ml := range req.MetricLimits {
			utility.Assert(ml.MetricId > 0, "Invalid metricId")
			utility.Assert(ml.MetricLimit > 0, "Invalid MetricLimit")
			me := query.GetMerchantMetric(ctx, ml.MetricId)
			utility.Assert(me != nil, "Metric not found")
			utility.Assert(me.Type == metric.MetricTypeLimitMetered, "Metric type invalid")
		}
	}

	if req.ProductName == nil || len(*req.ProductName) == 0 {
		req.ProductName = req.PlanName
	}
	if req.ProductDescription == nil || len(*req.ProductDescription) == 0 {
		req.ProductDescription = req.Description
	}

	if len(req.AddonIds) > 0 {
		var allAddonList []*entity.Plan
		err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, req.AddonIds).OmitEmpty().Scan(&allAddonList)
		for _, addonPlan := range allAddonList {
			utility.Assert(addonPlan.Type == consts.PlanTypeRecurringAddon, fmt.Sprintf("plan not recurring addon type, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("addon not published status, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Currency == one.Currency, fmt.Sprintf("addon currency not match plan's currency, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.ProductId == one.ProductId, fmt.Sprintf("addon product not match plan's product, id:%d", addonPlan.Id))
		}
	}

	if len(req.OnetimeAddonIds) > 0 {
		var allAddonList []*entity.Plan
		err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, req.OnetimeAddonIds).OmitEmpty().Scan(&allAddonList)
		for _, addonPlan := range allAddonList {
			utility.Assert(addonPlan.Type == consts.PlanTypeOnetime, fmt.Sprintf("plan not onetime addon type, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("addon not published status, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Currency == one.Currency, fmt.Sprintf("addon currency not match plan's currency, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.ProductId == one.ProductId, fmt.Sprintf("addon product not match plan's product, id:%d", addonPlan.Id))
		}
	}

	var editCurrency *string = nil
	if req.Currency != nil {
		editCurrency = unibee.String(strings.ToUpper(*req.Currency))
	}
	var bindingAddonIds *string = nil
	if req.AddonIds != nil {
		bindingAddonIds = unibee.String(utility.IntListToString(req.AddonIds))
	}
	var bindingOnetimeAddonIds *string = nil
	if req.OnetimeAddonIds != nil {
		bindingOnetimeAddonIds = unibee.String(utility.IntListToString(req.OnetimeAddonIds))
	}

	if req.MetricMeteredCharge != nil {
		utility.AssertError(MetricPlanChargeValidation(*req.MetricMeteredCharge), "Usage-based Metered charges validation failed")
		metricPlanCharge.MetricMeteredCharge = *req.MetricMeteredCharge
	}
	if req.MetricRecurringCharge != nil {
		utility.AssertError(MetricPlanChargeValidation(*req.MetricRecurringCharge), "Usage-based Recurring charges validation failed")
		metricPlanCharge.MetricRecurringCharge = *req.MetricRecurringCharge
	}
	if len(req.MetricLimits) > 0 {
		metricPlanCharge.MetricLimits = req.MetricLimits
	}

	utility.Assert(req.TrialDemand == nil || *req.TrialDemand == "" || *req.TrialDemand == "paymentMethod", "Demand of trial should be paymentMethod or not")

	_, err = dao.Plan.Ctx(ctx).Data(g.Map{
		dao.Plan.Columns().ExternalPlanId:         req.ExternalPlanId,
		dao.Plan.Columns().PlanName:               req.PlanName,
		dao.Plan.Columns().InternalName:           req.InternalName,
		dao.Plan.Columns().Amount:                 req.Amount,
		dao.Plan.Columns().Currency:               editCurrency,
		dao.Plan.Columns().IntervalUnit:           req.IntervalUnit,
		dao.Plan.Columns().IntervalCount:          req.IntervalCount,
		dao.Plan.Columns().Description:            req.Description,
		dao.Plan.Columns().ImageUrl:               req.ImageUrl,
		dao.Plan.Columns().HomeUrl:                req.HomeUrl,
		dao.Plan.Columns().BindingAddonIds:        bindingAddonIds,
		dao.Plan.Columns().BindingOnetimeAddonIds: bindingOnetimeAddonIds,
		dao.Plan.Columns().GasPayer:               req.GasPayer,
		dao.Plan.Columns().IsDeleted:              0,
		dao.Plan.Columns().TrialDemand:            req.TrialDemand,
		dao.Plan.Columns().TrialDurationTime:      req.TrialDurationTime,
		dao.Plan.Columns().TrialAmount:            req.TrialAmount,
		dao.Plan.Columns().CancelAtTrialEnd:       req.CancelAtTrialEnd,
		dao.Plan.Columns().ProductId:              req.ProductId,
		dao.Plan.Columns().MetricCharge:           utility.MarshalToJsonString(metricPlanCharge),
	}).Where(dao.Plan.Columns().Id, req.PlanId).OmitNil().Update()
	if err != nil {
		return nil, gerror.Newf(`PlanEdit record insert failure %s`, err)
	}
	if req.Metadata != nil {
		_, _ = dao.Plan.Ctx(ctx).Data(g.Map{
			dao.Plan.Columns().MetaData: utility.MarshalToJsonString(req.Metadata),
		}).Where(dao.Plan.Columns().Id, req.PlanId).OmitNil().Update()
	}

	one = query.GetPlanById(ctx, req.PlanId)

	if len(req.MetricLimits) > 0 {
		err = metric.BulkMetricLimitPlanBindingReplace(ctx, one, req.MetricLimits)
		if err != nil {
			return nil, gerror.Newf(`BulkMetricLimitPlanBindingReplace %s`, err)
		}
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Plan(%v)", one.Id),
		Content:        "Edit",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         one.Id,
		DiscountCode:   "",
	}, err)
	if req.AddonIds != nil && len(req.AddonIds) > 0 {
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("AddonIds(%s)", utility.IntListToString(req.AddonIds)),
			Content:        "OverrideRecurringAddons",
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         one.Id,
			DiscountCode:   "",
		}, err)
	}
	if req.OnetimeAddonIds != nil && len(req.OnetimeAddonIds) > 0 {
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("OneTimeAddonIds(%s)", utility.IntListToString(req.OnetimeAddonIds)),
			Content:        "OverrideOneTimeAddons",
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         one.Id,
			DiscountCode:   "",
		}, err)
	}
	return one, nil
}

func PlanCopy(ctx context.Context, planId uint64) (one *entity.Plan, err error) {
	utility.Assert(planId > 0, "PlanId should > 0")
	one = query.GetPlanById(ctx, planId)
	utility.Assert(one != nil, fmt.Sprintf("plan not found, id:%d", planId))
	one = &entity.Plan{
		CompanyId:              one.CompanyId,
		MerchantId:             one.MerchantId,
		PlanName:               one.PlanName + "(Copy)",
		InternalName:           one.InternalName,
		Amount:                 one.Amount,
		Currency:               one.Currency,
		IntervalUnit:           one.IntervalUnit,
		IntervalCount:          one.IntervalCount,
		Type:                   one.Type,
		Description:            one.Description,
		ImageUrl:               one.ImageUrl,
		HomeUrl:                one.HomeUrl,
		BindingAddonIds:        one.BindingAddonIds,
		BindingOnetimeAddonIds: one.BindingOnetimeAddonIds,
		Status:                 consts.PlanStatusEditable,
		CreateTime:             gtime.Now().Timestamp(),
		MetaData:               one.MetaData,
		PublishStatus:          consts.PlanPublishStatusUnPublished,
		GasPayer:               one.GasPayer,
		TrialDurationTime:      one.TrialDurationTime,
		TrialAmount:            one.TrialAmount,
		TrialDemand:            one.TrialDemand,
		CancelAtTrialEnd:       one.CancelAtTrialEnd,
		ProductId:              one.ProductId,
		MetricCharge:           one.MetricCharge,
	}
	result, err := dao.Plan.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`PlanCopy record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Plan(%v)", planId),
		Content:        fmt.Sprintf("CopyTo(%v)", one.Id),
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         planId,
		DiscountCode:   "",
	}, err)
	return one, nil
}

func PlanDelete(ctx context.Context, planId uint64) (one *entity.Plan, err error) {
	utility.Assert(planId > 0, "planId invalid")
	one = query.GetPlanById(ctx, planId)
	utility.Assert(one != nil, fmt.Sprintf("plan not found, id:%d", planId))
	utility.Assert(one.Status == consts.PlanStatusEditable, fmt.Sprintf("plan is not in edit status, id:%d", planId))
	_, err = dao.Plan.Ctx(ctx).Data(g.Map{
		dao.Plan.Columns().IsDeleted: gtime.Now().Timestamp(),
		dao.Plan.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Plan.Columns().Id, one.Id).Update()
	if err != nil {
		return nil, err
	}
	one.IsDeleted = int(gtime.Now().Timestamp())
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Plan(%v)", one.Id),
		Content:        "Delete",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         one.Id,
		DiscountCode:   "",
	}, err)
	return one, nil
}

func PlanArchive(ctx context.Context, planId uint64, hardArchive bool) (one *entity.Plan, err error) {
	utility.Assert(planId > 0, "planId invalid")
	one = query.GetPlanById(ctx, planId)
	utility.Assert(one != nil, fmt.Sprintf("plan not found, id:%d", planId))
	if !hardArchive {
		_, err = dao.Plan.Ctx(ctx).Data(g.Map{
			dao.Plan.Columns().Status:    consts.PlanStatusSoftArchive,
			dao.Plan.Columns().IsDeleted: -1,
			dao.Plan.Columns().GmtModify: gtime.Now(),
		}).Where(dao.Plan.Columns().Id, one.Id).Update()
		if err != nil {
			return nil, err
		}
		one.Status = consts.PlanStatusSoftArchive
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("Plan(%v)", one.Id),
			Content:        "SoftArchive",
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         one.Id,
			DiscountCode:   "",
		}, err)
	} else {
		_, err = dao.Plan.Ctx(ctx).Data(g.Map{
			dao.Plan.Columns().Status:            consts.PlanStatusHardArchive,
			dao.Plan.Columns().DisableAutoCharge: 1,
			dao.Plan.Columns().IsDeleted:         -2,
			dao.Plan.Columns().GmtModify:         gtime.Now(),
		}).Where(dao.Plan.Columns().Id, one.Id).Update()
		if err != nil {
			return nil, err
		}
		one.Status = consts.PlanStatusHardArchive
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("Plan(%v)", one.Id),
			Content:        "HardArchive",
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         one.Id,
			DiscountCode:   "",
		}, err)
	}
	return one, nil
}

func HardDeletePlan(ctx context.Context, planId uint64) error {
	_, err := dao.Plan.Ctx(ctx).Where(dao.Plan.Columns().Id, planId).Delete()
	if err != nil {
		return err
	}
	return nil
}

func PlanAddonsBinding(ctx context.Context, req *plan.AddonsBindingReq) (one *entity.Plan, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.Action >= 0 && req.Action <= 2, "action should 0-2")
	utility.Assert(req.PlanId > 0, "PlanId should > 0")
	one = query.GetPlanById(ctx, req.PlanId)
	utility.Assert(one != nil, fmt.Sprintf("plan not found, id:%d", req.PlanId))
	utility.Assert(one.Type == consts.PlanTypeMain, fmt.Sprintf("plan not type main, id:%d", req.PlanId))

	var addonIdsList []int64
	if len(one.BindingAddonIds) > 0 {
		//init
		strList := strings.Split(one.BindingAddonIds, ",")

		for _, s := range strList {
			num, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fmt.Println("Internal Error converting string to int:", err)
				return nil, err
			}
			addonIdsList = append(addonIdsList, num)
		}
	}
	var oneTimeAddonIdsList []int64
	if len(one.BindingOnetimeAddonIds) > 0 {
		//init
		strList := strings.Split(one.BindingOnetimeAddonIds, ",")

		for _, s := range strList {
			num, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fmt.Println("Internal Error converting string to int:", err)
				return nil, err
			}
			oneTimeAddonIdsList = append(oneTimeAddonIdsList, num)
		}
	}
	//addonIds type verify
	{
		var allAddonList []*entity.Plan
		err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, req.AddonIds).Scan(&allAddonList)
		for _, addonPlan := range allAddonList {
			utility.Assert(addonPlan.Type == consts.PlanTypeRecurringAddon, fmt.Sprintf("plan not addon recurring type, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("addon not published status, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Currency == one.Currency, fmt.Sprintf("addon currency not match plan's currency, id:%d", addonPlan.Id))
			//addon interval verify
			utility.Assert(addonPlan.IntervalUnit == one.IntervalUnit && addonPlan.IntervalCount == one.IntervalCount, fmt.Sprintf("addon not match plan's recycle interval, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.ProductId == one.ProductId, fmt.Sprintf("addon product not match plan's product, id:%d", addonPlan.Id))
		}
	}

	//onetime addonIds type verify
	{
		var allAddonList []*entity.Plan
		err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, req.OnetimeAddonIds).Scan(&allAddonList)
		for _, addonPlan := range allAddonList {
			utility.Assert(addonPlan.Type == consts.PlanTypeOnetime, fmt.Sprintf("plan not addon onetime type, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("addon not published status, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.Currency == one.Currency, fmt.Sprintf("addon currency not match plan's currency, id:%d", addonPlan.Id))
			utility.Assert(addonPlan.ProductId == one.ProductId, fmt.Sprintf("addon product not match plan's product, id:%d", addonPlan.Id))
		}
	}

	if req.Action == 0 {
		//replace
		addonIdsList = req.AddonIds
		oneTimeAddonIdsList = req.OnetimeAddonIds
	} else if req.Action == 1 {
		//add
		utility.Assert(len(req.AddonIds) > 0, "action add, addon ids is empty")
		addonIdsList = utility.MergeInt64Arrays(addonIdsList, req.AddonIds)
		oneTimeAddonIdsList = utility.MergeInt64Arrays(oneTimeAddonIdsList, req.OnetimeAddonIds)
	} else if req.Action == 2 {
		//delete
		utility.Assert(len(req.AddonIds) > 0, "action delete, addon ids is empty")
		addonIdsList = utility.RemoveInt64Arrays(addonIdsList, req.AddonIds)
		oneTimeAddonIdsList = utility.RemoveInt64Arrays(oneTimeAddonIdsList, req.OnetimeAddonIds)
	}

	utility.Assert(len(addonIdsList) <= 10, "binding addon too much, should <= 10")
	utility.Assert(len(oneTimeAddonIdsList) <= 10, "binding addon too much, should <= 10")

	one.BindingAddonIds = utility.IntListToString(addonIdsList)
	one.BindingOnetimeAddonIds = utility.IntListToString(oneTimeAddonIdsList)
	_, err = dao.Plan.Ctx(ctx).Data(g.Map{
		dao.Plan.Columns().BindingAddonIds:        one.BindingAddonIds,
		dao.Plan.Columns().BindingOnetimeAddonIds: one.BindingOnetimeAddonIds,
		dao.Plan.Columns().IsDeleted:              0,
		dao.Plan.Columns().GmtModify:              gtime.Now(),
	}).Where(dao.Plan.Columns().Id, one.Id).Update()

	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Plan(%v)", one.Id),
		Content:        fmt.Sprintf("BindingAddon-(%s)-(%s)", utility.MarshalToJsonString(req.AddonIds), utility.MarshalToJsonString(req.OnetimeAddonIds)),
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         one.Id,
		DiscountCode:   "",
	}, err)
	if err != nil {
		return nil, err
	}
	return one, nil
}
