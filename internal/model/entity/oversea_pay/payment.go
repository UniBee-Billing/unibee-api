// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Payment is the golang structure for table payment.
type Payment struct {
	Id                     int64       `json:"id"                     description:"id"`                                                                     // id
	CompanyId              int64       `json:"companyId"              description:"company id"`                                                             // company id
	MerchantId             int64       `json:"merchantId"             description:"merchant id"`                                                            // merchant id
	OpenApiId              int64       `json:"openApiId"              description:"open api id"`                                                            // open api id
	UserId                 int64       `json:"userId"                 description:"user_id"`                                                                // user_id
	SubscriptionId         string      `json:"subscriptionId"         description:"subscription id"`                                                        // subscription id
	GmtCreate              *gtime.Time `json:"gmtCreate"              description:"create time"`                                                            // create time
	BizType                int         `json:"bizType"                description:"biz_type 1-single payment, 3-subscription"`                              // biz_type 1-single payment, 3-subscription
	BizId                  string      `json:"bizId"                  description:"biz_id"`                                                                 // biz_id
	Currency               string      `json:"currency"               description:"currency，“SGD” “MYR” “PHP” “IDR” “THB”"`                                 // currency，“SGD” “MYR” “PHP” “IDR” “THB”
	PaymentId              string      `json:"paymentId"              description:"payment id"`                                                             // payment id
	TotalAmount            int64       `json:"totalAmount"            description:"total amount"`                                                           // total amount
	PaymentAmount          int64       `json:"paymentAmount"          description:"payment_amount"`                                                         // payment_amount
	BalanceAmount          int64       `json:"balanceAmount"          description:"balance_amount"`                                                         // balance_amount
	RefundAmount           int64       `json:"refundAmount"           description:"total refund amount"`                                                    // total refund amount
	Status                 int         `json:"status"                 description:"status  10-pending，20-success，30-failure, 40-cancel"`                    // status  10-pending，20-success，30-failure, 40-cancel
	TerminalIp             string      `json:"terminalIp"             description:"client ip"`                                                              // client ip
	CountryCode            string      `json:"countryCode"            description:"country code"`                                                           // country code
	AuthorizeStatus        int         `json:"authorizeStatus"        description:"authorize status，0-waiting authorize，1-authorized，2-authorized_request"` // authorize status，0-waiting authorize，1-authorized，2-authorized_request
	AuthorizeReason        string      `json:"authorizeReason"        description:""`                                                                       //
	ChannelId              int64       `json:"channelId"              description:"channel_id"`                                                             // channel_id
	ChannelPaymentIntentId string      `json:"channelPaymentIntentId" description:"channel_payment_intent_id"`                                              // channel_payment_intent_id
	ChannelPaymentId       string      `json:"channelPaymentId"       description:"channel_payment_id"`                                                     // channel_payment_id
	CaptureDelayHours      int         `json:"captureDelayHours"      description:"capture_delay_hours"`                                                    // capture_delay_hours
	CreateTime             *gtime.Time `json:"createTime"             description:"create time"`                                                            // create time
	CancelTime             *gtime.Time `json:"cancelTime"             description:"cancel time"`                                                            // cancel time
	PaidTime               *gtime.Time `json:"paidTime"               description:"paid time"`                                                              // paid time
	InvoiceId              string      `json:"invoiceId"              description:"invoice id"`                                                             // invoice id
	GmtModify              *gtime.Time `json:"gmtModify"              description:"update time"`                                                            // update time
	AppId                  string      `json:"appId"                  description:"app id"`                                                                 // app id
	ReturnUrl              string      `json:"returnUrl"              description:"return url"`                                                             // return url
	ChannelEdition         string      `json:"channelEdition"         description:"channel edition"`                                                        // channel edition
	HidePaymentMethods     string      `json:"hidePaymentMethods"     description:"hide_payment_methods"`                                                   // hide_payment_methods
	Verify                 string      `json:"verify"                 description:"codeVerify"`                                                             // codeVerify
	Code                   string      `json:"code"                   description:""`                                                                       //
	Token                  string      `json:"token"                  description:""`                                                                       //
	AdditionalData         string      `json:"additionalData"         description:"addtional data (json)"`                                                  // addtional data (json)
	Automatic              int         `json:"automatic"              description:""`                                                                       //
	FailureReason          string      `json:"failureReason"          description:""`                                                                       //
	BillingReason          string      `json:"billingReason"          description:""`                                                                       //
	Link                   string      `json:"link"                   description:""`                                                                       //
	PaymentData            string      `json:"paymentData"            description:"payment create data (json)"`                                             // payment create data (json)
	UniqueId               string      `json:"uniqueId"               description:"unique id"`                                                              // unique id
	BalanceStart           int64       `json:"balanceStart"           description:"balance_start"`                                                          // balance_start
	BalanceEnd             int64       `json:"balanceEnd"             description:"balance_end"`                                                            // balance_end
	InvoiceData            string      `json:"invoiceData"            description:""`                                                                       //
	ChannelPaymentMethod   string      `json:"channelPaymentMethod"   description:""`                                                                       //
}
