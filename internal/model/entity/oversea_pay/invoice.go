// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Invoice is the golang structure for table invoice.
type Invoice struct {
	Id                             uint64      `json:"id"                             description:""`                                                                       //
	MerchantId                     uint64      `json:"merchantId"                     description:"merchant_id"`                                                            // merchant_id
	UserId                         uint64      `json:"userId"                         description:"userId"`                                                                 // userId
	SubscriptionId                 string      `json:"subscriptionId"                 description:"subscription_id"`                                                        // subscription_id
	InvoiceId                      string      `json:"invoiceId"                      description:"invoice_id"`                                                             // invoice_id
	InvoiceName                    string      `json:"invoiceName"                    description:"invoice name"`                                                           // invoice name
	UniqueId                       string      `json:"uniqueId"                       description:"unique_id"`                                                              // unique_id
	GmtCreate                      *gtime.Time `json:"gmtCreate"                      description:"create time"`                                                            // create time
	GmtModify                      *gtime.Time `json:"gmtModify"                      description:"update time"`                                                            // update time
	TotalAmount                    int64       `json:"totalAmount"                    description:"total amount, cent"`                                                     // total amount, cent
	TaxAmount                      int64       `json:"taxAmount"                      description:"tax amount,cent"`                                                        // tax amount,cent
	SubscriptionAmount             int64       `json:"subscriptionAmount"             description:"sub amount,cent"`                                                        // sub amount,cent
	Currency                       string      `json:"currency"                       description:"currency"`                                                               // currency
	Lines                          string      `json:"lines"                          description:"lines( json)"`                                                           // lines( json)
	GatewayId                      uint64      `json:"gatewayId"                      description:"gateway_id"`                                                             // gateway_id
	Status                         int         `json:"status"                         description:"status，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled"` // status，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
	SendStatus                     int         `json:"sendStatus"                     description:"email send status，0-No | 1- YES| 2-Unnecessary"`                         // email send status，0-No | 1- YES| 2-Unnecessary
	SendEmail                      string      `json:"sendEmail"                      description:"email"`                                                                  // email
	SendPdf                        string      `json:"sendPdf"                        description:"pdf link"`                                                               // pdf link
	IsDeleted                      int         `json:"isDeleted"                      description:"0-UnDeleted，1-Deleted"`                                                  // 0-UnDeleted，1-Deleted
	Link                           string      `json:"link"                           description:"invoice link"`                                                           // invoice link
	PaymentLink                    string      `json:"paymentLink"                    description:"invoice payment link"`                                                   // invoice payment link
	GatewayStatus                  string      `json:"gatewayStatus"                  description:""`                                                                       //
	GatewayInvoiceId               string      `json:"gatewayInvoiceId"               description:""`                                                                       //
	GatewayPaymentId               string      `json:"gatewayPaymentId"               description:""`                                                                       //
	GatewayInvoicePdf              string      `json:"gatewayInvoicePdf"              description:""`                                                                       //
	TaxScale                       int64       `json:"taxScale"                       description:"Tax scale，1000 = 10%"`                                                   // Tax scale，1000 = 10%
	SendNote                       string      `json:"sendNote"                       description:"send_note"`                                                              // send_note
	SendTerms                      string      `json:"sendTerms"                      description:"send_terms"`                                                             // send_terms
	TotalAmountExcludingTax        int64       `json:"totalAmountExcludingTax"        description:""`                                                                       //
	SubscriptionAmountExcludingTax int64       `json:"subscriptionAmountExcludingTax" description:""`                                                                       //
	PeriodStart                    int64       `json:"periodStart"                    description:"period_start, utc time"`                                                 // period_start, utc time
	PeriodEnd                      int64       `json:"periodEnd"                      description:"period_end utc time"`                                                    // period_end utc time
	PeriodStartTime                *gtime.Time `json:"periodStartTime"                description:""`                                                                       //
	PeriodEndTime                  *gtime.Time `json:"periodEndTime"                  description:""`                                                                       //
	PaymentId                      string      `json:"paymentId"                      description:"paymentId"`                                                              // paymentId
	RefundId                       string      `json:"refundId"                       description:"refundId"`                                                               // refundId
	Data                           string      `json:"data"                           description:"data (json)"`                                                            // data (json)
	BizType                        int         `json:"bizType"                        description:"biz type from payment 1-single payment, 3-subscription"`                 // biz type from payment 1-single payment, 3-subscription
	CreateTime                     int64       `json:"createTime"                     description:"create utc time"`                                                        // create utc time
	CryptoAmount                   int64       `json:"cryptoAmount"                   description:"crypto_amount, cent"`                                                    // crypto_amount, cent
	CryptoCurrency                 string      `json:"cryptoCurrency"                 description:"crypto_currency"`                                                        // crypto_currency
	FinishTime                     int64       `json:"finishTime"                     description:"finish utc time"`                                                        // finish utc time
	DayUtilDue                     int64       `json:"dayUtilDue"                     description:"day util due after finish"`                                              // day util due after finish
}
