package user

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consts"
	"unibee/internal/logic/auth"
	"unibee/internal/logic/batch/export"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type TaskUserBillingExport struct {
}

func (t TaskUserBillingExport) TaskName() string {
	return "UserBillingExport"
}

func (t TaskUserBillingExport) Header() interface{} {
	return ExportUserBillingEntity{}
}

func (t TaskUserBillingExport) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
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
	req := &auth.UserListInternalReq{
		MerchantId: task.MerchantId,
		//UserId:        0,
		//Email:         "",
		//FirstName:     "",
		//LastName:      "",
		//Status:        nil,
		//DeleteInclude: false,
		//SortField:     "",
		//SortType:      "",
		//CreateTimeStart: 0,
		//CreateTimeEnd:   0,
		Page:  page,
		Count: count,
	}
	if payload != nil {
		if value, ok := payload["userId"].(float64); ok {
			req.UserId = int64(value)
		}
		if value, ok := payload["email"].(string); ok {
			req.Email = value
		}
		if value, ok := payload["firstName"].(string); ok {
			req.FirstName = value
		}
		if value, ok := payload["lastName"].(string); ok {
			req.LastName = value
		}
		if value, ok := payload["subscriptionId"].(string); ok {
			req.SubscriptionId = value
		}
		if value, ok := payload["status"].([]interface{}); ok {
			req.Status = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["subStatus"].([]interface{}); ok {
			req.SubStatus = export.JsonArrayTypeConvert(ctx, value)
		}
		//if value, ok := payload["deleteInclude"].(bool); ok {
		//	req.DeleteInclude = value
		//}
		if value, ok := payload["sortField"].(string); ok {
			req.SortField = value
		}
		if value, ok := payload["sortType"].(string); ok {
			req.SortType = value
		}
		if value, ok := payload["createTimeStart"].(float64); ok {
			req.CreateTimeStart = int64(value)
		}
		if value, ok := payload["createTimeEnd"].(float64); ok {
			req.CreateTimeEnd = int64(value)
		}
	}
	req.SkipTotal = true
	result, _ := auth.UserList(ctx, req)
	if result != nil && result.UserAccounts != nil {
		for _, one := range result.UserAccounts {
			var userGateway = ""
			if one.Gateway != nil {
				userGateway = one.Gateway.GatewayName
			}
			var recurringAmount = ""
			var currency = ""
			var billingPeriod = ""
			sub := query.GetSubscriptionBySubscriptionId(ctx, one.SubscriptionId)
			if sub != nil {
				currency = sub.Currency
				recurringAmount = utility.ConvertCentToDollarStr(one.RecurringAmount, sub.Currency)
				plan := query.GetPlanById(ctx, sub.PlanId)
				if plan != nil {
					billingPeriod = plan.IntervalUnit
				}
			}
			mainList = append(mainList, &ExportUserBillingEntity{
				Id:                 fmt.Sprintf("%v", one.Id),
				FirstName:          one.FirstName,
				LastName:           one.LastName,
				FullName:           fmt.Sprintf("%s %s", one.FirstName, one.LastName),
				Email:              one.Email,
				MerchantName:       merchant.Name,
				Phone:              one.Phone,
				Address:            one.Address,
				VATNumber:          one.VATNumber,
				CountryCode:        one.CountryCode,
				CountryName:        one.CountryName,
				SubscriptionName:   one.SubscriptionName,
				SubscriptionId:     one.SubscriptionId,
				PlanId:             fmt.Sprintf("%v", one.PlanId),
				RecurringAmount:    recurringAmount,
				Currency:           currency,
				BillingPeriod:      billingPeriod,
				SubscriptionStatus: consts.SubStatusToEnum(one.SubscriptionStatus).Description(),
				CreateTime:         gtime.NewFromTimeStamp(one.CreateTime),
				ExternalUserId:     one.ExternalUserId,
				Status:             consts.UserStatusToEnum(one.Status).Description(),
				TaxPercentage:      utility.ConvertTaxPercentageToPercentageString(one.TaxPercentage),
				Type:               consts.UserTypeToEnum(one.Type).Description(),
				Gateway:            userGateway,
				City:               one.City,
				ZipCode:            one.ZipCode,
			})
		}
	}
	return mainList, nil
}

type ExportUserBillingEntity struct {
	Id                 string      `json:"Id"                 `
	ExternalUserId     string      `json:"ExternalUserId"     `
	FirstName          string      `json:"FirstName"          `
	LastName           string      `json:"LastName"           `
	FullName           string      `json:"FullName"           `
	Email              string      `json:"Email"              `
	MerchantName       string      `json:"MerchantName"       `
	Phone              string      `json:"Phone"              `
	Address            string      `json:"Address"            `
	VATNumber          string      `json:"VATNumber"          `
	CountryCode        string      `json:"CountryCode"        `
	CountryName        string      `json:"CountryName"        `
	PlanId             string      `json:"PlanId"     `
	Currency           string      `json:"Currency"     `
	BillingPeriod      string      `json:"BillingPeriod"     `
	SubscriptionName   string      `json:"SubscriptionName"   `
	SubscriptionId     string      `json:"SubscriptionId"     `
	SubscriptionStatus string      `json:"SubscriptionStatus" `
	RecurringAmount    string      `json:"RecurringAmount" `
	CreateTime         *gtime.Time `json:"CreateTime"       layout:"2006-01-02 15:04:05"  `
	Status             string      `json:"Status"             `
	TaxPercentage      string      `json:"TaxPercentage"      `
	Type               string      `json:"Type"               `
	Gateway            string      `json:"Gateway"            `
	City               string      `json:"City"               `
	ZipCode            string      `json:"ZipCode"            `
}
