package detail

import (
	"context"
	"unibee/api/bean"
	"unibee/internal/consts"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

type PaymentTimelineDetail struct {
	Id                    uint64        `json:"id"             description:""`                                //
	MerchantId            uint64        `json:"merchantId"     description:"merchant id"`                     // merchant id
	UserId                uint64        `json:"userId"         description:"userId"`                          // userId
	SubscriptionId        string        `json:"subscriptionId" description:"subscription id"`                 // subscription id
	InvoiceId             string        `json:"invoiceId"      description:"invoice id"`                      // invoice id
	Currency              string        `json:"currency"       description:"currency"`                        // currency
	TotalAmount           int64         `json:"totalAmount"    description:"total amount"`                    // total amount
	GatewayId             uint64        `json:"gatewayId"      description:"gateway id"`                      // gateway id
	TransactionId         string        `json:"transactionId"      description:"TransactionId"`               // TransactionId
	PaymentId             string        `json:"paymentId"      description:"PaymentId"`                       // PaymentId
	Status                int           `json:"status"         description:"0-pending, 1-success, 2-failure"` // 0-pending, 1-success, 2-failure
	TimelineType          int           `json:"timelineType"   description:"0-pay, 1-refund"`                 // 0-pay, 1-refund
	CreateTime            int64         `json:"createTime"     description:"create utc time"`                 // create utc time
	RefundId              string        `json:"refundId"       description:"refund id"`                       // refund id
	FullRefund            int           `json:"fullRefund"     description:"0-no, 1-yes"`                     // 0-no, 1-yes
	Payment               *bean.Payment `json:"payment" dc:"Payment"`
	Refund                *bean.Refund  `json:"refund" dc:"Refund"`
	ExternalTransactionId string        `json:"externalTransactionId"      description:"ExternalTransactionId"` // ExternalTransactionId
}

func ConvertPaymentTimeline(ctx context.Context, one *entity.PaymentTimeline) *PaymentTimelineDetail {
	if one == nil {
		return nil
	}
	var payment = bean.SimplifyPayment(query.GetPaymentByPaymentId(ctx, one.PaymentId))
	if payment == nil {
		return nil
	}
	var refund = bean.SimplifyRefund(query.GetRefundByRefundId(ctx, one.RefundId))
	var transactionId = one.PaymentId
	var externalTransactionId = payment.GatewayPaymentId
	if one.TimelineType == consts.TimelineTypeRefund && refund != nil {
		transactionId = one.RefundId
		externalTransactionId = refund.GatewayRefundId
	}
	return &PaymentTimelineDetail{
		Id:                    one.Id,
		TransactionId:         transactionId,
		ExternalTransactionId: externalTransactionId,
		MerchantId:            one.MerchantId,
		UserId:                one.UserId,
		SubscriptionId:        one.SubscriptionId,
		InvoiceId:             one.InvoiceId,
		Currency:              one.Currency,
		TotalAmount:           one.TotalAmount,
		GatewayId:             one.GatewayId,
		PaymentId:             one.PaymentId,
		Status:                one.Status,
		TimelineType:          one.TimelineType,
		CreateTime:            one.CreateTime,
		RefundId:              one.RefundId,
		FullRefund:            one.FullRefund,
		Payment:               payment,
		Refund:                refund,
	}
}
