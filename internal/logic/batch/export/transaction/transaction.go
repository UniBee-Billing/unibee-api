package transaction

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/logic/batch/export"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"
)

type TaskTransactionExport struct {
}

func (t TaskTransactionExport) TaskName() string {
	return "TransactionExport"
}

func (t TaskTransactionExport) Header() interface{} {
	return ExportTransactionEntity{}
}

func (t TaskTransactionExport) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
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
		if value, ok := payload["userId"].(float64); ok {
			req.UserId = uint64(value)
		}
		if value, ok := payload["sortField"].(string); ok {
			req.SortField = value
		}
		if value, ok := payload["sortType"].(string); ok {
			req.SortType = value
		}
		if value, ok := payload["currency"].(string); ok {
			req.Currency = value
		}
		if value, ok := payload["createTimeStart"].(float64); ok {
			req.CreateTimeStart = int64(value)
		}
		if value, ok := payload["createTimeEnd"].(float64); ok {
			req.CreateTimeEnd = int64(value)
		}
		if value, ok := payload["amountStart"].(float64); ok {
			req.AmountStart = unibee.Int64(int64(value))
		}
		if value, ok := payload["amountEnd"].(float64); ok {
			req.AmountEnd = unibee.Int64(int64(value))
		}
		if value, ok := payload["status"].([]interface{}); ok {
			req.Status = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["timelineTypes"].([]interface{}); ok {
			req.TimelineTypes = export.JsonArrayTypeConvert(ctx, value)
		}
		if value, ok := payload["gatewayIds"].([]interface{}); ok {
			req.GatewayIds = export.JsonArrayTypeConvertUint64(ctx, value)
		}
	}
	req.SkipTotal = true
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
			} else {
				user = &entity.UserAccount{}
			}

			mainList = append(mainList, &ExportTransactionEntity{
				TransactionId:         one.TransactionId,
				UserId:                fmt.Sprintf("%v", user.Id),
				ExternalUserId:        fmt.Sprintf("%v", user.ExternalUserId),
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
	TransactionId         string      `json:"TransactionId"       comment:""`
	ExternalTransactionId string      `json:"externalTransactionId"   comment:""`
	UserId                string      `json:"UserId"              comment:""`
	ExternalUserId        string      `json:"ExternalUserId"      comment:""`
	FirstName             string      `json:"FirstName"           comment:""`
	LastName              string      `json:"LastName"            comment:""`
	Email                 string      `json:"Email"               comment:""`
	MerchantName          string      `json:"MerchantName"        comment:""`
	SubscriptionId        string      `json:"SubscriptionId"  comment:""`
	InvoiceId             string      `json:"InvoiceId"       comment:""`
	Currency              string      `json:"Currency"        comment:""`
	TotalAmount           string      `json:"TotalAmount"     comment:""`
	Gateway               string      `json:"Gateway"       comment:""`
	PaymentId             string      `json:"PaymentId"      comment:""`
	Status                string      `json:"Status"          comment:""`
	Type                  string      `json:"Type"    comment:""`
	CreateTime            *gtime.Time `json:"CreateTime"      layout:"2006-01-02 15:04:05" comment:""`
	RefundId              string      `json:"RefundId"       comment:""`
	FullRefund            string      `json:"FullRefund"      comment:""`
}
