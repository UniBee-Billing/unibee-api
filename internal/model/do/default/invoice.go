// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Invoice is the golang structure of table invoice for DAO operations like Where/Data.
type Invoice struct {
	g.Meta                         `orm:"table:invoice, do:true"`
	Id                             interface{} //
	MetaData                       interface{} // meta_data(json)
	BizType                        interface{} // biz type from payment 1-single payment, 3-subscription
	MerchantId                     interface{} // merchant_id
	UserId                         interface{} // userId
	SubscriptionId                 interface{} // subscription_id
	InvoiceId                      interface{} // invoice_id
	InvoiceName                    interface{} // invoice name
	UniqueId                       interface{} // unique_id
	GmtCreate                      *gtime.Time // create time
	GmtModify                      *gtime.Time // update time
	TotalAmount                    interface{} // total amount, cent
	TaxAmount                      interface{} // tax amount,cent
	SubscriptionAmount             interface{} // sub amount,cent
	Currency                       interface{} // currency
	Lines                          interface{} // lines( json)
	PaymentId                      interface{} // paymentId
	GatewayId                      interface{} // gateway_id
	Status                         interface{} // status，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
	SendStatus                     interface{} // email send status，0-No | 1- YES| 2-Unnecessary
	SendEmail                      interface{} // email
	SendPdf                        interface{} // pdf link
	IsDeleted                      interface{} // 0-UnDeleted，1-Deleted
	Link                           interface{} // invoice link
	PaymentLink                    interface{} // invoice payment link
	GatewayStatus                  interface{} //
	GatewayInvoiceId               interface{} //
	GatewayPaymentId               interface{} //
	GatewayInvoicePdf              interface{} //
	TaxPercentage                  interface{} // TaxPercentage，1000 = 10%
	SendNote                       interface{} // send_note
	SendTerms                      interface{} // send_terms
	TotalAmountExcludingTax        interface{} //
	SubscriptionAmountExcludingTax interface{} //
	PeriodStart                    interface{} // period_start, utc time
	PeriodEnd                      interface{} // period_end utc time
	TrialEnd                       interface{} // trial_end, utc time
	PeriodStartTime                *gtime.Time //
	PeriodEndTime                  *gtime.Time //
	RefundId                       interface{} // refundId
	Data                           interface{} // data (json)
	CreateTime                     interface{} // create utc time
	CryptoAmount                   interface{} // crypto_amount, cent
	CryptoCurrency                 interface{} // crypto_currency
	FinishTime                     interface{} // utc time of enter process
	DayUtilDue                     interface{} // day util due after process
	LastTrackTime                  interface{} // last process invoice track time
	DiscountCode                   interface{} // discount_code
	DiscountAmount                 interface{} // discount amount, cent
	CountryCode                    interface{} //
	ProductName                    interface{} // product name
	GasPayer                       interface{} // gas_payer
	GatewayPaymentMethod           interface{} // gateway_payment_method
	BillingCycleAnchor             interface{} // billing_cycle_anchor
	CreateFrom                     interface{} // create from
	VatNumber                      interface{} //
	PromoCreditDiscountAmount      interface{} // promo credit discount amount
	PartialCreditPaidAmount        interface{} // partial credit paid amount
	MetricCharge                   interface{} // invoice metric charge data
}
