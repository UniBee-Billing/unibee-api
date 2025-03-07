package bean

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type GatewayCurrencyExchange struct {
	FromCurrency   string  `json:"from_currency" description:"the currency of gateway exchange from"`
	ToCurrency     string  `json:"to_currency" description:"the currency of gateway exchange to"`
	ExchangeRate   float64 `json:"exchange_rate"  description:"the exchange rate of gateway, set to 0 if using https://app.exchangerate-api.com/ instead of fixed exchange rate"`
	ExchangeAmount int64   `json:"exchange_amount" description:"the exchange amount of gateway"`
}

var GatewayCurrencyExchangeKey = "GatewayCurrencyExchange"

type Payment struct {
	PaymentId               string                   `json:"paymentId"              description:"payment id"`                                                             // payment id
	MerchantId              uint64                   `json:"merchantId"             description:"merchant id"`                                                            // merchant id
	UserId                  uint64                   `json:"userId"                 description:"user_id"`                                                                // user_id
	SubscriptionId          string                   `json:"subscriptionId"         description:"subscription id"`                                                        // subscription id
	ExternalPaymentId       string                   `json:"externalPaymentId"      description:"external_payment_id"`                                                    // external_payment_id
	Currency                string                   `json:"currency"               description:"currency，“SGD” “MYR” “PHP” “IDR” “THB”"`                                 // currency，“SGD” “MYR” “PHP” “IDR” “THB”
	TotalAmount             int64                    `json:"totalAmount"            description:"total amount"`                                                           // total amount
	PaymentAmount           int64                    `json:"paymentAmount"          description:"payment_amount"`                                                         // payment_amount
	BalanceAmount           int64                    `json:"balanceAmount"          description:"balance_amount"`                                                         // balance_amount
	RefundAmount            int64                    `json:"refundAmount"           description:"total refund amount"`                                                    // total refund amount
	Status                  int                      `json:"status"                 description:"status  10-pending，20-success，30-failure, 40-cancel"`                    // status  10-pending，20-success，30-failure, 40-cancel
	CountryCode             string                   `json:"countryCode"            description:"country code"`                                                           // country code
	AuthorizeStatus         int                      `json:"authorizeStatus"        description:"authorize status，0-waiting authorize，1-authorized，2-authorized_request"` // authorize status，0-waiting authorize，1-authorized，2-authorized_request
	AuthorizeReason         string                   `json:"authorizeReason"        description:""`                                                                       //
	GatewayId               uint64                   `json:"gatewayId"              description:"gateway_id"`                                                             // gateway_id
	GatewayPaymentId        string                   `json:"gatewayPaymentId"       description:"gateway_payment_id"`                                                     // gateway_payment_id
	GatewayPaymentType      string                   `json:"gatewayPaymentType"       description:"gateway_payment_type"`
	CreateTime              int64                    `json:"createTime"             description:"create time, utc time"` // create time, utc time
	CancelTime              int64                    `json:"cancelTime"             description:"cancel time, utc time"` // cancel time, utc time
	PaidTime                int64                    `json:"paidTime"               description:"paid time, utc time"`   // paid time, utc time
	InvoiceId               string                   `json:"invoiceId"              description:"invoice id"`            // invoice id
	ReturnUrl               string                   `json:"returnUrl"              description:"return url"`            // return url
	Automatic               int                      `json:"automatic"              description:""`                      //
	FailureReason           string                   `json:"failureReason"          description:""`                      //
	BillingReason           string                   `json:"billingReason"          description:""`                      //
	Link                    string                   `json:"link"                   description:""`
	Metadata                map[string]interface{}   `json:"metadata"               description:""`
	GasPayer                string                   `json:"gasPayer"               description:"who pay the gas, merchant|user"` // who pay the gas, merchant|user
	AutoCharge              bool                     `json:"autoCharge"             description:""`
	GatewayCurrencyExchange *GatewayCurrencyExchange `json:"gatewayCurrencyExchange"  description:"gateway currency exchange"`
}

func SimplifyPayment(one *entity.Payment) *Payment {
	if one == nil {
		return nil
	}
	var metadata = make(map[string]interface{})
	var gatewayExchange *GatewayCurrencyExchange
	if len(one.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("SimplifyPayment Unmarshal Metadata error:%s", err.Error())
		}
	}
	var lastErr = one.LastError
	if len(lastErr) == 0 {
		lastErr = one.AuthorizeReason
	}
	autoCharge := false
	if len(one.GatewayPaymentMethod) > 0 {
		autoCharge = true
	}
	if s, ok := metadata[GatewayCurrencyExchangeKey]; ok {
		if value, ok2 := s.(string); ok2 {
			if len(value) > 0 {
				_ = utility.UnmarshalFromJsonString(value, &gatewayExchange)
			}
		}
	}
	if gatewayExchange != nil && len(gatewayExchange.FromCurrency) > 0 {
		gatewayExchange.ExchangeAmount = utility.ExchangeCurrencyConvert(one.TotalAmount, gatewayExchange.FromCurrency, gatewayExchange.ToCurrency, gatewayExchange.ExchangeRate)
	}
	return &Payment{
		PaymentId:               one.PaymentId,
		MerchantId:              one.MerchantId,
		UserId:                  one.UserId,
		SubscriptionId:          one.SubscriptionId,
		ExternalPaymentId:       one.ExternalPaymentId,
		Currency:                one.Currency,
		TotalAmount:             one.TotalAmount,
		PaymentAmount:           one.PaymentAmount,
		BalanceAmount:           one.BalanceAmount,
		RefundAmount:            one.RefundAmount,
		Status:                  one.Status,
		CountryCode:             one.CountryCode,
		AuthorizeStatus:         one.AuthorizeStatus,
		AuthorizeReason:         lastErr,
		GatewayId:               one.GatewayId,
		GatewayPaymentId:        one.GatewayPaymentId,
		GatewayPaymentType:      one.GatewayEdition,
		CreateTime:              one.CreateTime,
		CancelTime:              one.CancelTime,
		PaidTime:                one.PaidTime,
		InvoiceId:               one.InvoiceId,
		ReturnUrl:               one.ReturnUrl,
		Automatic:               one.Automatic,
		FailureReason:           one.FailureReason,
		BillingReason:           one.BillingReason,
		Link:                    one.Link,
		Metadata:                metadata,
		GasPayer:                one.GasPayer,
		AutoCharge:              autoCharge,
		GatewayCurrencyExchange: gatewayExchange,
	}
}

func SimplifyPaymentList(ones []*entity.Payment) (list []*Payment) {
	if len(ones) == 0 {
		return make([]*Payment, 0)
	}
	for _, one := range ones {
		list = append(list, SimplifyPayment(one))
	}
	return list
}
