package transaction

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type TaskTransaction struct {
}

func (t TaskTransaction) TaskName() string {
	return "TransactionExport"
}

func (t TaskTransaction) Header() interface{} {
	return ExportTransactionEntity{}
}

func (t TaskTransaction) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
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
	req := &service.PaymentTimelineListInternalReq{
		MerchantId: task.MerchantId,
		//UserId:     0,
		//CreateTimeStart: 0,
		//CreateTimeEnd:   0,
		//SortField: "",
		//SortType: "",
		Page:  page,
		Count: count,
	}
	if payload != nil {
		if value, ok := payload["userId"].(uint64); ok {
			req.UserId = value
		}
		if value, ok := payload["sortField"].(string); ok {
			req.SortField = value
		}
		if value, ok := payload["sortType"].(string); ok {
			req.SortType = value
		}
		if value, ok := payload["createTimeStart"].(int64); ok {
			req.CreateTimeStart = value
		}
		if value, ok := payload["createTimeEnd"].(int64); ok {
			req.CreateTimeEnd = value
		}
	}
	result, _ := service.PaymentTimeLineList(ctx, req)
	if result != nil && result.PaymentTimelines != nil {
		for _, one := range result.PaymentTimelines {
			if one == nil {
				continue
			}
			var gateway = ""
			gatewayEntity := query.GetGatewayById(ctx, one.GatewayId)
			if gatewayEntity != nil {
				gateway = gatewayEntity.GatewayName
			}
			user := query.GetUserAccountById(ctx, one.UserId)
			var transactionType = "payment"
			var fullRefund = "No"
			if one.TimelineType == 1 {
				transactionType = "refund"
				if one.FullRefund == 1 {
					fullRefund = "Yes"
				}
			}
			var status = "Pending"
			if one.Status == 1 {
				status = "Success"
			} else if one.Status == 2 {
				status = "Failure"
			}
			var firstName = ""
			var lastName = ""
			var email = ""
			if user != nil {
				firstName = user.FirstName
				lastName = user.LastName
				email = user.Email
			}

			mainList = append(mainList, &ExportTransactionEntity{
				TransactionId:         one.TransactionId,
				FirstName:             firstName,
				LastName:              lastName,
				Email:                 email,
				MerchantName:          merchant.Name,
				SubscriptionId:        one.SubscriptionId,
				InvoiceId:             one.InvoiceId,
				Currency:              one.Currency,
				TotalAmount:           utility.ConvertCentToDollarStr(one.TotalAmount, one.Currency),
				Gateway:               gateway,
				PaymentId:             one.PaymentId,
				Status:                status,
				Type:                  transactionType,
				CreateTime:            gtime.NewFromTimeStamp(one.CreateTime),
				RefundId:              one.RefundId,
				FullRefund:            fullRefund,
				ExternalTransactionId: one.ExternalTransactionId,
			})
		}
	}
	return mainList, nil
}

type ExportTransactionEntity struct {
	TransactionId         string      `json:"TransactionId"      `
	FirstName             string      `json:"FirstName"          `
	LastName              string      `json:"LastName"           `
	Email                 string      `json:"Email"              `
	MerchantName          string      `json:"MerchantName"       `
	SubscriptionId        string      `json:"SubscriptionId" `
	InvoiceId             string      `json:"InvoiceId"      `
	Currency              string      `json:"Currency"       `
	TotalAmount           string      `json:"TotalAmount"    `
	Gateway               string      `json:"Gateway"      `
	PaymentId             string      `json:"PaymentId"     `
	Status                string      `json:"Status"         `
	Type                  string      `json:"Type"   `
	CreateTime            *gtime.Time `json:"CreateTime"      layout:"2006-01-02 15:04:05"`
	RefundId              string      `json:"RefundId"      `
	FullRefund            string      `json:"FullRefund"     `
	ExternalTransactionId string      `json:"externalTransactionId"  `
}
