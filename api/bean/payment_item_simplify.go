package bean

import (
	entity "unibee/internal/model/entity/default"
)

type PaymentItemSimplify struct {
	Id             uint64 `json:"id"             description:""`                                           //
	BizType        int    `json:"bizType"        description:"biz_type 1-onetime payment, 3-subscription"` // biz_type 1-onetime payment, 3-subscription
	MerchantId     uint64 `json:"merchantId"     description:"merchant id"`                                // merchant id
	UserId         uint64 `json:"userId"         description:"userId"`                                     // userId
	SubscriptionId string `json:"subscriptionId" description:"subscription id"`                            // subscription id
	InvoiceId      string `json:"invoiceId"      description:"invoice id"`                                 // invoice id
	UniqueId       string `json:"uniqueId"       description:"unique id"`                                  // unique id
	Currency       string `json:"currency"       description:"currency"`                                   // currency
	Amount         int64  `json:"amount"         description:"amount"`                                     // amount
	UnitAmount     int64  `json:"unitAmount"     description:"unit_amount"`                                // unit_amount
	Quantity       int64  `json:"quantity"       description:"quantity"`                                   // quantity
	PaymentId      string `json:"paymentId"      description:"PaymentId"`                                  // PaymentId
	Status         int    `json:"status"         description:"0-pending, 1-success, 2-failure"`            // 0-pending, 1-success, 2-failure
	CreateTime     int64  `json:"createTime"     description:"create utc time"`                            // create utc time
	Description    string `json:"description"    description:"description"`                                // description
	Name           string `json:"name"           description:"name"`                                       // name
}

func SimplifyPaymentItemTimeline(one *entity.PaymentItem) *PaymentItemSimplify {
	if one == nil {
		return nil
	}
	return &PaymentItemSimplify{
		Id:             one.Id,
		BizType:        one.BizType,
		MerchantId:     one.MerchantId,
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		Name:           one.Name,
		Description:    one.Description,
		Currency:       one.Currency,
		Amount:         one.Amount,
		UnitAmount:     one.UnitAmount,
		Quantity:       one.Quantity,
		PaymentId:      one.PaymentId,
		Status:         one.Status,
		CreateTime:     one.CreateTime,
	}
}
