package bean

import entity "unibee/internal/model/entity/oversea_pay"

type PaymentTimelineSimplify struct {
	Id             uint64 `json:"id"             description:""`                                //
	MerchantId     uint64 `json:"merchantId"     description:"merchant id"`                     // merchant id
	UserId         uint64 `json:"userId"         description:"userId"`                          // userId
	SubscriptionId string `json:"subscriptionId" description:"subscription id"`                 // subscription id
	InvoiceId      string `json:"invoiceId"      description:"invoice id"`                      // invoice id
	Currency       string `json:"currency"       description:"currency"`                        // currency
	TotalAmount    int64  `json:"totalAmount"    description:"total amount"`                    // total amount
	GatewayId      uint64 `json:"gatewayId"      description:"gateway id"`                      // gateway id
	PaymentId      string `json:"paymentId"      description:"PaymentId"`                       // PaymentId
	Status         int    `json:"status"         description:"0-pending, 1-success, 2-failure"` // 0-pending, 1-success, 2-failure
	TimelineType   int    `json:"timelineType"   description:"0-pay, 1-refund"`                 // 0-pay, 1-refund
	CreateTime     int64  `json:"createTime"     description:"create utc time"`                 // create utc time
	RefundId       string `json:"refundId"       description:"refund id"`                       // refund id
	FullRefund     int    `json:"fullRefund"     description:"0-no, 1-yes"`                     // 0-no, 1-yes

}

func SimplifyPaymentTimeline(one *entity.PaymentTimeline) *PaymentTimelineSimplify {
	if one == nil {
		return nil
	}
	return &PaymentTimelineSimplify{
		Id:             one.Id,
		MerchantId:     one.MerchantId,
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		Currency:       one.Currency,
		TotalAmount:    one.TotalAmount,
		GatewayId:      one.GatewayId,
		PaymentId:      one.PaymentId,
		Status:         one.Status,
		TimelineType:   one.TimelineType,
		CreateTime:     one.CreateTime,
		FullRefund:     one.FullRefund,
		RefundId:       one.RefundId,
	}
}