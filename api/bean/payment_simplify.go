package bean

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	entity "unibee/internal/model/entity/oversea_pay"
)

type PaymentSimplify struct {
	PaymentId         string            `json:"paymentId"              description:"payment id"`                                                             // payment id
	MerchantId        uint64            `json:"merchantId"             description:"merchant id"`                                                            // merchant id
	UserId            int64             `json:"userId"                 description:"user_id"`                                                                // user_id
	SubscriptionId    string            `json:"subscriptionId"         description:"subscription id"`                                                        // subscription id
	ExternalPaymentId string            `json:"externalPaymentId"      description:"external_payment_id"`                                                    // external_payment_id
	Currency          string            `json:"currency"               description:"currency，“SGD” “MYR” “PHP” “IDR” “THB”"`                                 // currency，“SGD” “MYR” “PHP” “IDR” “THB”
	TotalAmount       int64             `json:"totalAmount"            description:"total amount"`                                                           // total amount
	PaymentAmount     int64             `json:"paymentAmount"          description:"payment_amount"`                                                         // payment_amount
	BalanceAmount     int64             `json:"balanceAmount"          description:"balance_amount"`                                                         // balance_amount
	RefundAmount      int64             `json:"refundAmount"           description:"total refund amount"`                                                    // total refund amount
	Status            int               `json:"status"                 description:"status  10-pending，20-success，30-failure, 40-cancel"`                    // status  10-pending，20-success，30-failure, 40-cancel
	CountryCode       string            `json:"countryCode"            description:"country code"`                                                           // country code
	AuthorizeStatus   int               `json:"authorizeStatus"        description:"authorize status，0-waiting authorize，1-authorized，2-authorized_request"` // authorize status，0-waiting authorize，1-authorized，2-authorized_request
	AuthorizeReason   string            `json:"authorizeReason"        description:""`                                                                       //
	GatewayId         uint64            `json:"gatewayId"              description:"gateway_id"`                                                             // gateway_id
	GatewayPaymentId  string            `json:"gatewayPaymentId"       description:"gateway_payment_id"`                                                     // gateway_payment_id
	CreateTime        int64             `json:"createTime"             description:"create time, utc time"`                                                  // create time, utc time
	CancelTime        int64             `json:"cancelTime"             description:"cancel time, utc time"`                                                  // cancel time, utc time
	PaidTime          int64             `json:"paidTime"               description:"paid time, utc time"`                                                    // paid time, utc time
	InvoiceId         string            `json:"invoiceId"              description:"invoice id"`                                                             // invoice id
	ReturnUrl         string            `json:"returnUrl"              description:"return url"`                                                             // return url
	Automatic         int               `json:"automatic"              description:""`                                                                       //
	FailureReason     string            `json:"failureReason"          description:""`                                                                       //
	BillingReason     string            `json:"billingReason"          description:""`                                                                       //
	Link              string            `json:"link"                   description:""`
	Metadata          map[string]string `json:"metadata" description:""`
	GasPayer          string            `json:"gasPayer"                  description:"who pay the gas, merchant|user"` // who pay the gas, merchant|user
}

func SimplifyPayment(one *entity.Payment) *PaymentSimplify {
	if one == nil {
		return nil
	}
	var metadata = make(map[string]string)
	if len(one.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("SimplifyPayment Unmarshal Metadata error:%s", err.Error())
		}
	}
	return &PaymentSimplify{
		PaymentId:         one.PaymentId,
		MerchantId:        one.MerchantId,
		UserId:            one.UserId,
		SubscriptionId:    one.SubscriptionId,
		ExternalPaymentId: one.ExternalPaymentId,
		Currency:          one.Currency,
		TotalAmount:       one.TotalAmount,
		PaymentAmount:     one.PaymentAmount,
		BalanceAmount:     one.BalanceAmount,
		RefundAmount:      one.RefundAmount,
		Status:            one.Status,
		CountryCode:       one.CountryCode,
		AuthorizeStatus:   one.AuthorizeStatus,
		AuthorizeReason:   one.AuthorizeReason,
		GatewayId:         one.GatewayId,
		GatewayPaymentId:  one.GatewayPaymentId,
		CreateTime:        one.CreateTime,
		CancelTime:        one.CancelTime,
		PaidTime:          one.PaidTime,
		InvoiceId:         one.InvoiceId,
		ReturnUrl:         one.ReturnUrl,
		Automatic:         one.Automatic,
		FailureReason:     one.FailureReason,
		BillingReason:     one.BillingReason,
		Link:              one.Link,
		Metadata:          metadata,
		GasPayer:          one.GasPayer,
	}
}

func SimplifyPaymentList(ones []*entity.Payment) (list []*PaymentSimplify) {
	if len(ones) == 0 {
		return make([]*PaymentSimplify, 0)
	}
	for _, one := range ones {
		list = append(list, SimplifyPayment(one))
	}
	return list
}
