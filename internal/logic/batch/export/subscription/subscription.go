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
	entity "unibee/internal/model/entity/default"
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
			req.PlanIds = export.JsonArrayTypeConvertUint64(ctx, value)
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
			var stripeUserId = ""
			var stripePaymentMethod = ""
			var paypalVaultId = ""
			if one.Gateway != nil {
				subGateway = one.Gateway.GatewayName
				if one.Gateway.GatewayType == consts.GatewayTypeCard {
					gatewayUser := query.GetGatewayUser(ctx, one.Subscription.UserId, one.Gateway.Id)
					if gatewayUser != nil {
						stripeUserId = gatewayUser.GatewayUserId
						stripePaymentMethod = one.Subscription.DefaultPaymentMethodId
					}
				} else if one.Gateway.GatewayType == consts.GatewayTypePaypal {
					paypalVaultId = one.Subscription.DefaultPaymentMethodId
				}
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
				one.User = &bean.UserAccount{}
			}
			if one.Plan == nil {
				one.Plan = &bean.Plan{}
			}
			var productName = ""
			if one.Plan.ProductId > 0 {
				product := query.GetProductById(ctx, uint64(one.Plan.ProductId))
				if product != nil {
					productName = product.ProductName
				}
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
				ProductId:              fmt.Sprintf("%v", one.Plan.ProductId),
				ProductName:            productName,
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
				StripeUserId:           stripeUserId,
				StripePaymentMethod:    stripePaymentMethod,
				PaypalVaultId:          paypalVaultId,
			})
		}
	}
	return mainList, nil
}

type ExportSubscriptionEntity struct {
	SubscriptionId         string      `json:"SubscriptionId"     comment:""`
	ExternalSubscriptionId string      `json:"ExternalSubscriptionId"     comment:""`
	UserId                 string      `json:"UserId"             comment:""`
	ExternalUserId         string      `json:"ExternalUserId"     comment:""`
	ProductId              string      `json:"ProductId"             comment:""`
	ProductName            string      `json:"ProductName"             comment:""`
	PlanId                 string      `json:"PlanId"             comment:""`
	ExternalPlanId         string      `json:"ExternalPlanId"     comment:""`
	FirstName              string      `json:"FirstName"          comment:""`
	LastName               string      `json:"LastName"           comment:""`
	Email                  string      `json:"Email"              comment:""`
	MerchantName           string      `json:"MerchantName"       comment:""`
	Amount                 string      `json:"Amount"             comment:""`
	Currency               string      `json:"Currency"           comment:""`
	PlanName               string      `json:"PlanName"           comment:""`
	Quantity               string      `json:"Quantity"           comment:""`
	Gateway                string      `json:"Gateway"            comment:""`
	Status                 string      `json:"Status"             comment:""`
	CancelAtPeriodEnd      string      `json:"CancelAtPeriodEnd"  comment:""`
	CurrentPeriodStart     *gtime.Time `json:"CurrentPeriodStart" layout:"2006-01-02 15:04:05" comment:""`
	CurrentPeriodEnd       *gtime.Time `json:"CurrentPeriodEnd"  layout:"2006-01-02 15:04:05" comment:""`
	BillingCycleAnchor     *gtime.Time `json:"BillingCycleAnchor" layout:"2006-01-02 15:04:05" comment:""`
	DunningTime            *gtime.Time `json:"DunningTime"        layout:"2006-01-02 15:04:05" comment:""`
	TrialEnd               *gtime.Time `json:"TrialEnd"          layout:"2006-01-02 15:04:05"  comment:""`
	FirstPaidTime          *gtime.Time `json:"FirstPaidTime"    layout:"2006-01-02 15:04:05"   comment:""`
	CancelReason           string      `json:"CancelReason"        comment:""`
	CountryCode            string      `json:"CountryCode"         comment:""`
	TaxPercentage          string      `json:"TaxPercentage"       comment:""`
	CreateTime             *gtime.Time `json:"CreateTime"      layout:"2006-01-02 15:04:05"    comment:""`
	StripeUserId           string      `json:"StripeUserId(Auto-Charge Required)"      comment:"The id of user get from stripe, required if stripe auto-charge needed"       `
	StripePaymentMethod    string      `json:"StripePaymentMethod(Auto-Charge Required)"     comment:"The payment method id which user attached, get from stripe, required if stripe auto-charge needed"    `
	PaypalVaultId          string      `json:"PaypalVaultId(Auto-Charge Required)"    comment:"The vault id of user get from paypal, required if paypal auto-charge needed"   `
}
