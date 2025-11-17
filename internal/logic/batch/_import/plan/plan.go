package user

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	currency2 "unibee/internal/logic/currency"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/plan"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type TaskPlanImport struct {
}

func (t TaskPlanImport) TemplateVersion() string {
	return "v1"
}

func (t TaskPlanImport) TaskName() string {
	return "PlanImport"
}

func (t TaskPlanImport) TemplateHeader() interface{} {
	return &ImportPlanEntity{
		ExternalPlanId:    "examplePlanId",
		PlanName:          "testPlan",
		InternalName:      "testPlan",
		Description:       "testPlan",
		PlanType:          "MainPlan",
		PlanAmount:        "10",
		Currency:          "EUR",
		IntervalUnit:      "month",
		IntervalCount:     "1",
		Metadata:          "",
		TrialAmount:       "0",
		TrialDurationTime: "0",
		TrialDemand:       "",
		CancelAtTrialEnd:  "false",
	}
}

func (t TaskPlanImport) ImportRow(ctx context.Context, task *entity.MerchantBatchTask, row map[string]string) (interface{}, error) {
	var err error
	target := &ImportPlanEntity{
		ExternalPlanId:    fmt.Sprintf("%s", row["ExternalPlanId"]),
		PlanName:          fmt.Sprintf("%s", row["PlanName"]),
		ProductName:       fmt.Sprintf("%s", row["ProductName"]),
		InternalName:      fmt.Sprintf("%s", row["InternalName"]),
		Description:       fmt.Sprintf("%s", row["Description"]),
		PlanType:          fmt.Sprintf("%s", row["PlanType"]),
		PlanAmount:        fmt.Sprintf("%s", row["PlanAmount"]),
		Currency:          fmt.Sprintf("%s", row["Currency"]),
		IntervalUnit:      fmt.Sprintf("%s", row["IntervalUnit"]),
		IntervalCount:     fmt.Sprintf("%s", row["IntervalCount"]),
		Metadata:          fmt.Sprintf("%s", row["Metadata"]),
		TrialAmount:       fmt.Sprintf("%s", row["TrialAmount"]),
		TrialDurationTime: fmt.Sprintf("%s", row["TrialDurationTime"]),
		TrialDemand:       fmt.Sprintf("%s", row["TrialDemand"]),
		CancelAtTrialEnd:  fmt.Sprintf("%s", row["CancelAtTrialEnd"]),
	}
	if len(target.PlanName) == 0 {
		return target, gerror.New("Error, PlanName is blank")
	}
	if len(target.ExternalPlanId) == 0 {
		return target, gerror.New("Error, ExternalPlanId is blank")
	}
	amountFloat, err := strconv.ParseFloat(target.PlanAmount, 64)
	if err != nil {
		return target, gerror.Newf("Invalid Amount,error:%s", err.Error())
	}
	amount := int64(amountFloat * 100)
	if amount <= 0 {
		return target, gerror.New("Invalid Amount, should greater than 0")
	}

	var trialAmount int64 = 0
	if len(target.TrialAmount) > 0 {
		trialAmountFloat, err := strconv.ParseFloat(target.TrialAmount, 64)
		if err != nil {
			return target, gerror.Newf("Invalid TrialAmount,error:%s", err.Error())
		}
		trialAmount = int64(trialAmountFloat * 100)
		if trialAmount < 0 {
			return target, gerror.New("Invalid TrialAmount, should greater than 0")
		}
	}

	if len(target.Currency) == 0 {
		return target, gerror.New("Error, Currency is blank")
	}
	currency := strings.TrimSpace(strings.ToUpper(target.Currency))
	if !currency2.IsCurrencySupport(currency) {
		return target, gerror.New("Error, invalid Currency")
	}
	if utility.IsNoCentCurrency(currency) {
		if amount%100 != 0 {
			return target, gerror.New("Error, this currency No decimals allowed，made it divisible by 100")
		}
	}
	product := query.GetProductByProductName(ctx, target.ProductName, task.MerchantId)
	if product == nil {
		return target, gerror.New("Error, product not found")
	}
	intervals := []string{"day", "month", "year", "week"}
	utility.Assert(utility.StringContainsElement(intervals, strings.ToLower(target.IntervalUnit)), "IntervalUnit Error， must one of day｜month｜year｜week\"")
	intervalCount, _ := strconv.ParseInt(target.IntervalCount, 10, 64)
	if intervalCount == 0 {
		intervalCount = 1
	}
	var trialDurationTime int64 = 0
	if len(target.TrialDurationTime) > 0 {
		trialDurationTime, _ = strconv.ParseInt(target.TrialDurationTime, 10, 64)
		if trialDurationTime < 0 {
			trialDurationTime = 0
		}
	}
	var cancelAtTrialEnd = 0
	if len(target.CancelAtTrialEnd) > 0 {
		cancelAtTrialEnd, _ = strconv.Atoi(target.CancelAtTrialEnd)
		if cancelAtTrialEnd < 0 || cancelAtTrialEnd > 1 {
			cancelAtTrialEnd = 0
		}
	}
	var metadata map[string]interface{}
	if len(target.Metadata) > 0 {
		_ = utility.UnmarshalFromJsonString(target.Metadata, &metadata)
	}
	one := query.GetPlanByExternalPlanId(ctx, task.MerchantId, target.ExternalPlanId)
	if one != nil {
		_, err = dao.Plan.Ctx(ctx).Data(g.Map{
			dao.Plan.Columns().MerchantId:        task.MerchantId,
			dao.Plan.Columns().PlanName:          target.PlanName,
			dao.Plan.Columns().InternalName:      target.InternalName,
			dao.Plan.Columns().Description:       target.Description,
			dao.Plan.Columns().MetaData:          metadata,
			dao.Plan.Columns().TrialAmount:       trialAmount,
			dao.Plan.Columns().TrialDurationTime: trialDurationTime,
			dao.Plan.Columns().TrialDemand:       target.TrialDemand,
			dao.Plan.Columns().CancelAtTrialEnd:  cancelAtTrialEnd,
			dao.Plan.Columns().GmtModify:         gtime.Now(),
		}).Where(dao.Plan.Columns().Id, one.Id).OmitEmpty().Update()

		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("Plan(%v)", one.Id),
			Content:        "ImportOverride",
			PlanId:         one.Id,
			SubscriptionId: "",
			InvoiceId:      "",
			DiscountCode:   "",
		}, err)

		if err == nil {
			err = gerror.New("override success")
		}
		return target, err
	}
	one, err = plan.PlanCreate(ctx, &plan.PlanInternalReq{
		ExternalPlanId:    target.ExternalPlanId,
		MerchantId:        task.MerchantId,
		PlanName:          target.PlanName,
		InternalName:      target.InternalName,
		Amount:            amount,
		Currency:          currency,
		IntervalUnit:      target.IntervalUnit,
		IntervalCount:     int(intervalCount),
		Description:       target.Description,
		Type:              int(consts.PlanTypeDescriptionToEnum(target.PlanType)),
		Metadata:          metadata,
		TrialAmount:       trialAmount,
		TrialDurationTime: trialDurationTime,
		TrialDemand:       target.TrialDemand,
		CancelAtTrialEnd:  cancelAtTrialEnd,
		ProductId:         int64(product.Id),
	})
	if err != nil {
		return target, err
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Plan(%v)", one.Id),
		Content:        "ImportNew",
		PlanId:         one.Id,
		SubscriptionId: "",
		InvoiceId:      "",
		DiscountCode:   "",
	}, err)
	return target, err
}

type ImportPlanEntity struct {
	ExternalPlanId    string `json:"ExternalPlanId"            comment:"external_user_id"`
	PlanName          string `json:"PlanName"                  comment:""`
	ProductName       string `json:"ProductName"               comment:""`
	InternalName      string `json:"InternalName"              comment:"PlanInternalName"`
	Description       string `json:"Description"               comment:"description"`
	PlanType          string `json:"PlanType"                  comment:""`
	PlanAmount        string `json:"PlanAmount"                comment:""`
	Currency          string `json:"Currency"                  comment:""`
	IntervalUnit      string `json:"IntervalUnit"              comment:"period unit,day|month|year|week"`
	IntervalCount     string `json:"IntervalCount"             comment:"period unit count"`
	Metadata          string `json:"Metadata"                  comment:""`
	TrialAmount       string `json:"TrialAmount"               comment:"price of trial period"`
	TrialDurationTime string `json:"TrialDurationTime"         comment:"duration of trial，seconds"`
	TrialDemand       string `json:"TrialDemand"               comment:""`
	CancelAtTrialEnd  string `json:"CancelAtTrialEnd"          comment:"whether cancel at subscription first trial end"`
}
