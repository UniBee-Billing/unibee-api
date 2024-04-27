package bean

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/internal/controller/link"
	entity "unibee/internal/model/entity/oversea_pay"
)

type InvoiceSimplify struct {
	Id                             uint64                 `json:"id"                             description:""`
	InvoiceId                      string                 `json:"invoiceId"`
	InvoiceName                    string                 `json:"invoiceName"`
	DiscountCode                   string                 `json:"discountCode"`
	OriginAmount                   int64                  `json:"originAmount"                `
	TotalAmount                    int64                  `json:"totalAmount"`
	DiscountAmount                 int64                  `json:"discountAmount"`
	TotalAmountExcludingTax        int64                  `json:"totalAmountExcludingTax"`
	Currency                       string                 `json:"currency"`
	TaxAmount                      int64                  `json:"taxAmount"`
	TaxPercentage                  int64                  `json:"taxPercentage"                  description:"TaxPercentage，1000 = 10%"`
	SubscriptionAmount             int64                  `json:"subscriptionAmount"`
	SubscriptionAmountExcludingTax int64                  `json:"subscriptionAmountExcludingTax"`
	Lines                          []*InvoiceItemSimplify `json:"lines"`
	PeriodEnd                      int64                  `json:"periodEnd"`
	PeriodStart                    int64                  `json:"periodStart"`
	FinishTime                     int64                  `json:"finishTime"`
	ProrationDate                  int64                  `json:"prorationDate"`
	ProrationScale                 int64                  `json:"prorationScale"`
	Link                           string                 `json:"link"                           description:"invoice link"` // invoice link
	PaymentLink                    string                 `json:"paymentLink"                    description:"invoice payment link"`
	Status                         int                    `json:"status"                         description:"status，1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled"` // status，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
	PaymentId                      string                 `json:"paymentId"                      description:"paymentId"`                                                     // paymentId
	BizType                        int                    `json:"bizType"                        description:"biz type from payment 1-onetime payment, 3-subscription"`       // biz type from payment 1-single payment, 3-subscription
	CryptoAmount                   int64                  `json:"cryptoAmount"                   description:"crypto_amount, cent"`                                           // crypto_amount, cent
	CryptoCurrency                 string                 `json:"cryptoCurrency"                 description:"crypto_currency"`
	SendStatus                     int                    `json:"sendStatus"                     description:"email send status，0-No | 1- YES| 2-Unnecessary"` // email send status，0-No | 1- YES| 2-Unnecessary
	DayUtilDue                     int64                  `json:"dayUtilDue"                     description:"day util due after finish"`                      // day util due after finish
}

type InvoiceItemSimplify struct {
	Currency               string        `json:"currency"`
	OriginAmount           int64         `json:"originAmount"`
	DiscountAmount         int64         `json:"discountAmount"`
	Amount                 int64         `json:"amount"`
	Tax                    int64         `json:"tax"`
	AmountExcludingTax     int64         `json:"amountExcludingTax"`
	TaxPercentage          int64         `json:"taxPercentage"                  description:"Tax Percentage，1000 = 10%"`
	UnitAmountExcludingTax int64         `json:"unitAmountExcludingTax"`
	Description            string        `json:"description"`
	Proration              bool          `json:"proration"`
	Quantity               int64         `json:"quantity"`
	PeriodEnd              int64         `json:"periodEnd"`
	PeriodStart            int64         `json:"periodStart"`
	Plan                   *PlanSimplify `json:"plan"`
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
		OriginAmount:                   one.TotalAmount + one.DiscountAmount,
		TotalAmount:                    one.TotalAmount,
		DiscountCode:                   one.DiscountCode,
		DiscountAmount:                 one.DiscountAmount,
		TotalAmountExcludingTax:        one.TotalAmountExcludingTax,
		Currency:                       one.Currency,
		TaxAmount:                      one.TaxAmount,
		SubscriptionAmount:             one.SubscriptionAmount,
		SubscriptionAmountExcludingTax: one.SubscriptionAmountExcludingTax,
		Lines:                          lines,
		PeriodEnd:                      one.PeriodEnd,
		PeriodStart:                    one.PeriodStart,
		FinishTime:                     one.FinishTime,
		Link:                           link.GetInvoiceLink(one.InvoiceId, one.SendTerms),
		PaymentLink:                    one.PaymentLink,
		Status:                         one.Status,
		PaymentId:                      one.PaymentId,
		BizType:                        one.BizType,
		CryptoCurrency:                 one.CryptoCurrency,
		CryptoAmount:                   one.CryptoAmount,
		SendStatus:                     one.SendStatus,
		DayUtilDue:                     one.DayUtilDue,
		TaxPercentage:                  one.TaxPercentage,
	}
}
