package bean

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	entity "unibee/internal/model/entity/oversea_pay"
)

type InvoiceSimplify struct {
	Id                             uint64                 `json:"id"                             description:""` //
	InvoiceId                      string                 `json:"invoiceId"`
	InvoiceName                    string                 `json:"invoiceName"`
	TotalAmount                    int64                  `json:"totalAmount"`
	TotalAmountExcludingTax        int64                  `json:"totalAmountExcludingTax"`
	Currency                       string                 `json:"currency"`
	TaxAmount                      int64                  `json:"taxAmount"`
	TaxScale                       int64                  `json:"taxScale"                  description:"Tax Scale，1000 = 10%"`
	SubscriptionAmount             int64                  `json:"subscriptionAmount"`
	SubscriptionAmountExcludingTax int64                  `json:"subscriptionAmountExcludingTax"`
	Lines                          []*InvoiceItemSimplify `json:"lines"`
	PeriodEnd                      int64                  `json:"periodEnd"`
	PeriodStart                    int64                  `json:"periodStart"`
	ProrationDate                  int64                  `json:"prorationDate"`
	ProrationScale                 int64                  `json:"prorationScale"`
	Link                           string                 `json:"link"                           description:"invoice link"` // invoice link
	PaymentLink                    string                 `json:"paymentLink"                    description:"invoice payment link"`
	Status                         int                    `json:"status"                         description:"status，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled"` // status，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
	PaymentId                      string                 `json:"paymentId"                      description:"paymentId"`                                                              // paymentId
	RefundId                       string                 `json:"refundId"                       description:"refundId"`                                                               // refundId
	BizType                        int                    `json:"bizType"                        description:"biz type from payment 1-single payment, 3-subscription"`                 // biz type from payment 1-single payment, 3-subscription
	CryptoAmount                   int64                  `json:"cryptoAmount"                   description:"crypto_amount, cent"`                                                    // crypto_amount, cent
	CryptoCurrency                 string                 `json:"cryptoCurrency"                 description:"crypto_currency"`
}

type InvoiceItemSimplify struct {
	Currency               string `json:"currency"`
	Amount                 int64  `json:"amount"`
	Tax                    int64  `json:"tax"`
	AmountExcludingTax     int64  `json:"amountExcludingTax"`
	TaxScale               int64  `json:"taxScale"                  description:"Tax Scale，1000 = 10%"`
	UnitAmountExcludingTax int64  `json:"unitAmountExcludingTax"`
	Description            string `json:"description"`
	Proration              bool   `json:"proration"`
	Quantity               int64  `json:"quantity"`
	PeriodEnd              int64  `json:"periodEnd"`
	PeriodStart            int64  `json:"periodStart"`
}

func UnmarshalFromJsonString(target string, one interface{}) error {
	if len(target) > 0 {
		return gjson.Unmarshal([]byte(target), &one)
	} else {
		return gerror.New("target is nil")
	}
}

func SimplifyInvoice(one *entity.Invoice) *InvoiceSimplify {
	if one == nil {
		return nil
	}
	var lines []*InvoiceItemSimplify
	err := UnmarshalFromJsonString(one.Lines, &lines)
	if err != nil {
		return nil
	}
	return &InvoiceSimplify{
		Id:                             one.Id,
		InvoiceId:                      one.InvoiceId,
		TotalAmount:                    one.TotalAmount,
		TotalAmountExcludingTax:        one.TotalAmountExcludingTax,
		Currency:                       one.Currency,
		TaxAmount:                      one.TaxAmount,
		SubscriptionAmount:             one.SubscriptionAmount,
		SubscriptionAmountExcludingTax: one.SubscriptionAmountExcludingTax,
		Lines:                          lines,
		PeriodEnd:                      one.PeriodEnd,
		PeriodStart:                    one.PeriodStart,
		Link:                           one.Link,
		PaymentLink:                    one.PaymentLink,
		Status:                         one.Status,
		PaymentId:                      one.PaymentId,
		RefundId:                       one.RefundId,
		BizType:                        one.BizType,
		CryptoCurrency:                 one.CryptoCurrency,
		CryptoAmount:                   one.CryptoAmount,
	}
}