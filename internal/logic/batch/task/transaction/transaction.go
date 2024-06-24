package transaction

import (
	"context"
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
	result, _ := service.PaymentTimeLineList(ctx, &service.PaymentTimelineListInternalReq{
		MerchantId: task.MerchantId,
		//UserId:     0,
		//CreateTimeStart: 0,
		//CreateTimeEnd:   0,
		Page:  page,
		Count: count,
	})
	if result != nil && result.PaymentTimelines != nil {
		for _, one := range result.PaymentTimelines {
			var gateway = ""
			gatewayEntity := query.GetGatewayById(ctx, one.GatewayId)
			if gatewayEntity != nil {
				gateway = gatewayEntity.GatewayName
			}
			user := query.GetUserAccountById(ctx, one.UserId)
			if user == nil {
				continue
			}
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

			mainList = append(mainList, &ExportTransactionEntity{
				TransactionId:         one.TransactionId,
				FirstName:             user.FirstName,
				LastName:              user.LastName,
				Email:                 user.Email,
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
