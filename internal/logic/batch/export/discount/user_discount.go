package discount

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/logic/discount"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type TaskUserDiscount struct {
}

func (t TaskUserDiscount) TaskName() string {
	return "UserDiscountExport"
}

func (t TaskUserDiscount) Header() interface{} {
	return ExportUserDiscountEntity{}
}

func (t TaskUserDiscount) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
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
	var id int64
	if value, ok := payload["id"].(int64); ok {
		id = value
	}
	if id <= 0 {
		return mainList, nil
	}
	req := &discount.UserDiscountListInternalReq{
		MerchantId: task.MerchantId,
		Id:         uint64(id),
		Page:       page,
		Count:      count,
	}
	if payload != nil {
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
	result, _ := discount.MerchantUserDiscountCodeList(ctx, req)
	if result != nil {
		for _, one := range result {
			var firstName = ""
			var lastName = ""
			var email = ""
			if one.User != nil {
				firstName = one.User.FirstName
				lastName = one.User.LastName
				email = one.User.Email
			}
			mainList = append(mainList, &ExportUserDiscountEntity{
				Id:             fmt.Sprintf("%v", one.Id),
				MerchantName:   merchant.Name,
				FirstName:      firstName,
				LastName:       lastName,
				Email:          email,
				PlanId:         fmt.Sprintf("%v", one.Plan.Id),
				PlanName:       one.Plan.PlanName,
				Code:           one.Code,
				SubscriptionId: one.SubscriptionId,
				PaymentId:      one.PaymentId,
				InvoiceId:      one.InvoiceId,
				CreateTime:     gtime.NewFromTimeStamp(one.CreateTime),
				ApplyAmount:    utility.ConvertCentToDollarStr(one.ApplyAmount, one.Currency),
				Currency:       one.Currency,
			})
		}
	}
	return mainList, nil
}

type ExportUserDiscountEntity struct {
	Id             string      `json:"Id"                 `
	FirstName      string      `json:"FirstName"          `
	LastName       string      `json:"LastName"           `
	Email          string      `json:"Email"              `
	MerchantName   string      `json:"MerchantName"       `
	PlanId         string      `json:"PlanId"             `
	PlanName       string      `json:"PlanName"           `
	Code           string      `json:"Code"           `
	SubscriptionId string      `json:"SubscriptionId" `
	PaymentId      string      `json:"PaymentId"      `
	InvoiceId      string      `json:"InvoiceId"      `
	CreateTime     *gtime.Time `json:"CreateTime"      layout:"2006-01-02 15:04:05"`
	ApplyAmount    string      `json:"ApplyAmount"    `
	Currency       string      `json:"Currency"       `
}
