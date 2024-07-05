package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/logic/batch/export"
	"unibee/internal/logic/subscription/service"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"
)

type TaskSubscriptionExport struct {
}

func (t TaskSubscriptionExport) TaskName() string {
	return "SubscriptionExport"
}

func (t TaskSubscriptionExport) Header() interface{} {
	return ExportSubscriptionEntity{}
}

func (t TaskSubscriptionExport) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
	var mainList = make([]interface{}, 0)
	if task == nil && task.MerchantId <= 0 {
		return mainList, nil
	}
	merchant := query.GetMerchantById(ctx, task.MerchantId)
	var payload map[string]interface{}
	err := utility.UnmarshalFromJsonString(task.Payload, &payload)
	if err != nil {
		g.Log().Errorf(ctx, "Download PageData error:%s", err.Error())
		return mainList, nil
	}
	req := &service.SubscriptionListInternalReq{
		MerchantId: task.MerchantId,
		//CreateTimeStart: 0,
		//CreateTimeEnd:   0,
		Page:  page,
		Count: count,
	}
	if payload != nil {
		if value, ok := payload["userId"].(float64); ok {
			req.UserId = int64(value)
		}
		if value, ok := payload["sortField"].(string); ok {
			req.SortField = value
		}
		if value, ok := payload["sortType"].(string); ok {
			req.SortType = value
		}
		if value, ok := payload["status"].([]interface{}); ok {
			req.Status = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["planIds"].([]interface{}); ok {
			req.PlanIds = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["currency"].(string); ok {
			req.Currency = value
		}
		if value, ok := payload["amountStart"].(float64); ok {
			req.AmountStart = unibee.Int64(int64(value))
		}
		if value, ok := payload["amountEnd"].(float64); ok {
			req.AmountEnd = unibee.Int64(int64(value))
		}
		if value, ok := payload["createTimeStart"].(float64); ok {
			req.CreateTimeStart = int64(value)
		}
		if value, ok := payload["createTimeEnd"].(float64); ok {
			req.CreateTimeEnd = int64(value)
		}
	}
	req.SkipTotal = true
	result, _ := service.SubscriptionList(ctx, req)
	if result != nil {
		for _, one := range result {
			var subGateway = ""
			if one.Gateway != nil {
				subGateway = one.Gateway.GatewayName
			}
			var canAtPeriodEnd = "No"
			if one.Subscription.CancelAtPeriodEnd == 1 {
				canAtPeriodEnd = "Yes"
			}
			var firstName = ""
			var lastName = ""
			var email = ""
			if one.User != nil {
				firstName = one.User.FirstName
				lastName = one.User.LastName
				email = one.User.Email
			} else {
				one.User = &bean.UserAccountSimplify{}
			}
			if one.Plan == nil {
				one.Plan = &bean.PlanSimplify{}
			}
			mainList = append(mainList, &ExportSubscriptionEntity{
				SubscriptionId:         one.Subscription.SubscriptionId,
				ExternalSubscriptionId: one.Subscription.ExternalSubscriptionId,
				UserId:                 fmt.Sprintf("%v", one.User.Id),
				ExternalUserId:         fmt.Sprintf("%v", one.User.ExternalUserId),
				FirstName:              firstName,
				LastName:               lastName,
				Email:                  email,
				MerchantName:           merchant.Name,
				Amount:                 utility.ConvertCentToDollarStr(one.Subscription.Amount, one.Subscription.Currency),
				Currency:               one.Subscription.Currency,
				PlanId:                 fmt.Sprintf("%v", one.Plan.Id),
				ExternalPlanId:         fmt.Sprintf("%v", one.Plan.ExternalPlanId),
				PlanName:               one.Plan.PlanName,
				Quantity:               fmt.Sprintf("%v", one.Subscription.Quantity),
				Gateway:                subGateway,
				Status:                 consts.SubStatusToEnum(one.Subscription.Status).Description(),
				CancelAtPeriodEnd:      canAtPeriodEnd,
				CurrentPeriodStart:     gtime.NewFromTimeStamp(one.Subscription.CurrentPeriodStart),
				CurrentPeriodEnd:       gtime.NewFromTimeStamp(one.Subscription.CurrentPeriodEnd),
				BillingCycleAnchor:     gtime.NewFromTimeStamp(one.Subscription.BillingCycleAnchor),
				DunningTime:            gtime.NewFromTimeStamp(one.Subscription.DunningTime),
				TrialEnd:               gtime.NewFromTimeStamp(one.Subscription.TrialEnd),
				FirstPaidTime:          gtime.NewFromTimeStamp(one.Subscription.FirstPaidTime),
				CancelReason:           one.Subscription.CancelReason,
				CountryCode:            one.Subscription.CountryCode,
				TaxPercentage:          utility.ConvertTaxPercentageToPercentageString(one.Subscription.TaxPercentage),
				CreateTime:             gtime.NewFromTimeStamp(one.Subscription.CreateTime),
			})
		}
	}
	return mainList, nil
}

type ExportSubscriptionEntity struct {
	SubscriptionId         string      `json:"SubscriptionId"     `
	ExternalSubscriptionId string      `json:"ExternalSubscriptionId"     `
	UserId                 string      `json:"UserId"             `
	ExternalUserId         string      `json:"ExternalUserId"     `
	PlanId                 string      `json:"PlanId"             `
	ExternalPlanId         string      `json:"ExternalPlanId"     `
	FirstName              string      `json:"FirstName"          `
	LastName               string      `json:"LastName"           `
	Email                  string      `json:"Email"              `
	MerchantName           string      `json:"MerchantName"       `
	Amount                 string      `json:"Amount"             `
	Currency               string      `json:"Currency"           `
	PlanName               string      `json:"PlanName"           `
	Quantity               string      `json:"Quantity"           `
	Gateway                string      `json:"Gateway"            `
	Status                 string      `json:"Status"             `
	CancelAtPeriodEnd      string      `json:"CancelAtPeriodEnd"  `
	CurrentPeriodStart     *gtime.Time `json:"CurrentPeriodStart" layout:"2006-01-02 15:04:05"`
	CurrentPeriodEnd       *gtime.Time `json:"CurrentPeriodEnd"  layout:"2006-01-02 15:04:05" `
	BillingCycleAnchor     *gtime.Time `json:"BillingCycleAnchor" layout:"2006-01-02 15:04:05"`
	DunningTime            *gtime.Time `json:"DunningTime"        layout:"2006-01-02 15:04:05"`
	TrialEnd               *gtime.Time `json:"TrialEnd"          layout:"2006-01-02 15:04:05" `
	FirstPaidTime          *gtime.Time `json:"FirstPaidTime"    layout:"2006-01-02 15:04:05"  `
	CancelReason           string      `json:"CancelReason"       `
	CountryCode            string      `json:"CountryCode"        `
	TaxPercentage          string      `json:"TaxPercentage"      `
	CreateTime             *gtime.Time `json:"CreateTime"      layout:"2006-01-02 15:04:05"   `
}
