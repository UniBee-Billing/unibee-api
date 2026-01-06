package transaction

import (
	"context"
	"fmt"
	"strings"
	"unibee/api/bean"
	"unibee/internal/logic/batch/export"
	"unibee/internal/logic/payment/service"
	preload2 "unibee/internal/logic/preload"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"

	dao "unibee/internal/dao/default"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type TaskTransactionV2Export struct {
}

func (t TaskTransactionV2Export) TaskName() string {
	return "TransactionExport"
}

func (t TaskTransactionV2Export) Header() interface{} {
	return ExportTransactionEntity{}
}

func (t TaskTransactionV2Export) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
	var mainList = make([]interface{}, 0)
	if task == nil || task.MerchantId <= 0 {
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
		Page:       page,
		Count:      count,
	}
	var timeZone int64 = 0
	timeZoneStr := fmt.Sprintf("UTC")
	if payload != nil {
		if value, ok := payload["timeZone"].(string); ok {
			zone, err := export.GetUTCOffsetFromTimeZone(value)
			if err == nil && zone > 0 {
				timeZoneStr = value
				timeZone = zone
			}
		}
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
			req.CreateTimeStart = int64(value) - timeZone
		}
		if value, ok := payload["createTimeEnd"].(float64); ok {
			req.CreateTimeEnd = int64(value) - timeZone
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

	// Get payment timelines directly from database to avoid N+1 queries in ConvertPaymentTimeline
	paymentTimelines := transactionList(ctx, req)
	if len(paymentTimelines) > 0 {
		// Preload all related data to avoid N+1 queries
		preload := preload2.TransactionListPreload(ctx, paymentTimelines)

		for _, one := range paymentTimelines {
			if one == nil {
				continue
			}

			// Use preloaded data instead of individual queries
			var gateway = ""
			if preload.Gateways[one.GatewayId] != nil {
				gateway = preload.Gateways[one.GatewayId].GatewayName
			}

			user := &entity.UserAccount{}
			if preload.Users[one.UserId] != nil {
				user = preload.Users[one.UserId]
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
			} else if one.Status == 3 {
				status = "Cancel"
			}

			var firstName = ""
			var lastName = ""
			var email = ""
			if user != nil {
				firstName = user.FirstName
				lastName = user.LastName
				email = user.Email
			}

			var exchangeCurrency = ""
			var exchangeAmount = ""
			var exchangeRate = ""

			// Get payment and refund data from preloaded data
			payment := bean.SimplifyPayment(preload.Payments[one.PaymentId])
			refund := bean.SimplifyRefund(preload.Refunds[one.RefundId])

			if refund == nil && payment != nil && payment.GatewayCurrencyExchange != nil {
				exchangeCurrency = payment.GatewayCurrencyExchange.ToCurrency
				exchangeAmount = utility.ConvertCentToDollarStr(payment.GatewayCurrencyExchange.ExchangeAmount, payment.GatewayCurrencyExchange.ToCurrency)
				exchangeRate = fmt.Sprintf("%.2f", payment.GatewayCurrencyExchange.ExchangeRate)
			} else if refund != nil && refund.GatewayCurrencyExchange != nil {
				exchangeCurrency = refund.GatewayCurrencyExchange.ToCurrency
				exchangeAmount = utility.ConvertCentToDollarStr(refund.GatewayCurrencyExchange.ExchangeAmount, refund.GatewayCurrencyExchange.ToCurrency)
				exchangeRate = fmt.Sprintf("%.2f", refund.GatewayCurrencyExchange.ExchangeRate)
			}

			mainList = append(mainList, &ExportTransactionEntity{
				TransactionId:         one.PaymentId, // Use PaymentId as TransactionId
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
				CreateTime:            gtime.NewFromTimeStamp(one.CreateTime + timeZone),
				RefundId:              one.RefundId,
				FullRefund:            fullRefund,
				ExternalTransactionId: getExternalTransactionId(payment, refund, one),
				TimeZone:              timeZoneStr,
				ExchangeAmount:        exchangeAmount,
				ExchangeCurrency:      exchangeCurrency,
				ExchangeRate:          exchangeRate,
			})
		}
	}
	return mainList, nil
}

// transactionList directly queries payment timeline data from database to avoid N+1 queries
func transactionList(ctx context.Context, req *service.PaymentTimelineListInternalReq) []*entity.PaymentTimeline {
	var mainList = make([]*entity.PaymentTimeline, 0)
	var total = 0
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}

	utility.Assert(req.MerchantId > 0, "merchantId not found")
	var sortKey = "gmt_create desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("merchant_id|gmt_create|gmt_modify|user_id", req.SortField), "sortField should one of merchant_id|gmt_create|gmt_modify|user_id")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " desc"
		} else {
			sortKey = req.SortField + " desc"
		}
	}

	q := dao.PaymentTimeline.Ctx(ctx).
		Where(dao.PaymentTimeline.Columns().MerchantId, req.MerchantId)

	if req.UserId > 0 {
		q = q.Where(dao.PaymentTimeline.Columns().UserId, req.UserId)
	}
	if req.CreateTimeStart > 0 {
		q = q.WhereGTE(dao.PaymentTimeline.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		q = q.WhereLTE(dao.PaymentTimeline.Columns().CreateTime, req.CreateTimeEnd)
	}
	if req.AmountStart != nil && req.AmountEnd != nil {
		utility.Assert(*req.AmountStart <= *req.AmountEnd, "amountStart should lower than amountEnd")
	}
	if req.AmountStart != nil {
		q = q.WhereGTE(dao.PaymentTimeline.Columns().TotalAmount, &req.AmountStart)
	}
	if req.AmountEnd != nil {
		q = q.WhereLTE(dao.PaymentTimeline.Columns().TotalAmount, &req.AmountEnd)
	}
	if len(req.Status) > 0 {
		q = q.WhereIn(dao.PaymentTimeline.Columns().Status, req.Status)
	}
	if len(req.TimelineTypes) > 0 {
		q = q.WhereIn(dao.PaymentTimeline.Columns().TimelineType, req.TimelineTypes)
	}
	if len(req.GatewayIds) > 0 {
		q = q.WhereIn(dao.PaymentTimeline.Columns().GatewayId, req.GatewayIds)
	}
	if len(req.Currency) > 0 {
		q = q.Where(dao.PaymentTimeline.Columns().Currency, strings.ToUpper(req.Currency))
	}

	q = q.Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty()

	var err error
	if req.SkipTotal {
		err = q.Scan(&mainList)
	} else {
		err = q.ScanAndCount(&mainList, &total, true)
	}
	if err != nil {
		return mainList
	}

	return mainList
}

// getExternalTransactionId gets the external transaction ID based on timeline type
func getExternalTransactionId(payment *bean.Payment, refund *bean.Refund, timeline *entity.PaymentTimeline) string {
	if timeline.TimelineType == 1 && refund != nil { // Refund type
		return refund.GatewayRefundId
	}
	if payment != nil {
		return payment.GatewayPaymentId
	}
	return ""
}
