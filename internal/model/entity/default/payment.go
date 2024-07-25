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
	MerchantId             uint64      `json:"merchantId"             description:"merchant id"`                                                            // merchant id
	LastError              string      `json:"lastError"              description:"last error"`                                                             // last error
	AuthorizeReason        string      `json:"authorizeReason"        description:""`                                                                       //
	Code                   string      `json:"code"                   description:""`                                                                       //
	GatewayPaymentMethod   string      `json:"gatewayPaymentMethod"   description:""`                                                                       //
	UserId                 uint64      `json:"userId"                 description:"user_id"`                                                                // user_id
	SubscriptionId         string      `json:"subscriptionId"         description:"subscription id"`                                                        // subscription id
	GmtCreate              *gtime.Time `json:"gmtCreate"              description:"create time"`                                                            // create time
	GmtModify              *gtime.Time `json:"gmtModify"              description:"update time"`                                                            // update time
	BizType                int         `json:"bizType"                description:"biz_type 1-onetime payment, 3-subscription"`                             // biz_type 1-onetime payment, 3-subscription
	ExternalPaymentId      string      `json:"externalPaymentId"      description:"external_payment_id"`                                                    // external_payment_id
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
	GatewayId              uint64      `json:"gatewayId"              description:"gateway_id"`                                                             // gateway_id
	GatewayPaymentIntentId string      `json:"gatewayPaymentIntentId" description:"gateway_payment_intent_id"`                                              // gateway_payment_intent_id
	GatewayPaymentId       string      `json:"gatewayPaymentId"       description:"gateway_payment_id"`                                                     // gateway_payment_id
	CaptureDelayHours      int         `json:"captureDelayHours"      description:"capture_delay_hours"`                                                    // capture_delay_hours
	CreateTime             int64       `json:"createTime"             description:"create time, utc time"`                                                  // create time, utc time
	CancelTime             int64       `json:"cancelTime"             description:"cancel time, utc time"`                                                  // cancel time, utc time
	PaidTime               int64       `json:"paidTime"               description:"paid time, utc time"`                                                    // paid time, utc time
	InvoiceId              string      `json:"invoiceId"              description:"invoice id"`                                                             // invoice id
	AppId                  string      `json:"appId"                  description:"app id"`                                                                 // app id
	ReturnUrl              string      `json:"returnUrl"              description:"return url"`                                                             // return url
	OpenApiId              int64       `json:"openApiId"              description:"open api id"`                                                            // open api id
	GatewayEdition         string      `json:"gatewayEdition"         description:"gateway edition"`                                                        // gateway edition
	HidePaymentMethods     string      `json:"hidePaymentMethods"     description:"hide_payment_methods"`                                                   // hide_payment_methods
	Verify                 string      `json:"verify"                 description:"codeVerify"`                                                             // codeVerify
	Token                  string      `json:"token"                  description:""`                                                                       //
	MetaData               string      `json:"metaData"               description:"meta_data (json)"`                                                       // meta_data (json)
	Automatic              int         `json:"automatic"              description:"0-no,1-yes"`                                                             // 0-no,1-yes
	FailureReason          string      `json:"failureReason"          description:""`                                                                       //
	BillingReason          string      `json:"billingReason"          description:""`                                                                       //
	Link                   string      `json:"link"                   description:""`                                                                       //
	PaymentData            string      `json:"paymentData"            description:"payment create data (json)"`                                             // payment create data (json)
	UniqueId               string      `json:"uniqueId"               description:"unique id"`                                                              // unique id
	BalanceStart           int64       `json:"balanceStart"           description:"balance_start, utc time"`                                                // balance_start, utc time
	BalanceEnd             int64       `json:"balanceEnd"             description:"balance_end, utc time"`                                                  // balance_end, utc time
	InvoiceData            string      `json:"invoiceData"            description:""`                                                                       //
	GasPayer               string      `json:"gasPayer"               description:"who pay the gas, merchant|user"`                                         // who pay the gas, merchant|user
	ExpireTime             int64       `json:"expireTime"             description:"expire time, utc time"`                                                  // expire time, utc time
	GatewayLink            string      `json:"gatewayLink"            description:""`                                                                       //
	CryptoAmount           int64       `json:"cryptoAmount"           description:"crypto_amount, cent"`                                                    // crypto_amount, cent
	CryptoCurrency         string      `json:"cryptoCurrency"         description:"crypto_currency"`                                                        // crypto_currency
}
