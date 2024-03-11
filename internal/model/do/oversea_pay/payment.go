// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Payment is the golang structure of table payment for DAO operations like Where/Data.
type Payment struct {
	g.Meta                 `orm:"table:payment, do:true"`
	Id                     interface{} // id
	CompanyId              interface{} // company id
	MerchantId             interface{} // merchant id
	OpenApiId              interface{} // open api id
	UserId                 interface{} // user_id
	SubscriptionId         interface{} // subscription id
	GmtCreate              *gtime.Time // create time
	GmtModify              *gtime.Time // update time
	BizType                interface{} // biz_type 1-single payment, 3-subscription
	ExternalPaymentId      interface{} // external_payment_id
	Currency               interface{} // currency，“SGD” “MYR” “PHP” “IDR” “THB”
	PaymentId              interface{} // payment id
	TotalAmount            interface{} // total amount
	PaymentAmount          interface{} // payment_amount
	BalanceAmount          interface{} // balance_amount
	RefundAmount           interface{} // total refund amount
	Status                 interface{} // status  10-pending，20-success，30-failure, 40-cancel
	TerminalIp             interface{} // client ip
	CountryCode            interface{} // country code
	AuthorizeStatus        interface{} // authorize status，0-waiting authorize，1-authorized，2-authorized_request
	AuthorizeReason        interface{} //
	GatewayId              interface{} // gateway_id
	GatewayPaymentIntentId interface{} // gateway_payment_intent_id
	GatewayPaymentId       interface{} // gateway_payment_id
	CaptureDelayHours      interface{} // capture_delay_hours
	CreateTime             interface{} // create time, utc time
	CancelTime             interface{} // cancel time, utc time
	PaidTime               interface{} // paid time, utc time
	InvoiceId              interface{} // invoice id
	AppId                  interface{} // app id
	ReturnUrl              interface{} // return url
	GatewayEdition         interface{} // gateway edition
	HidePaymentMethods     interface{} // hide_payment_methods
	Verify                 interface{} // codeVerify
	Code                   interface{} //
	Token                  interface{} //
	MetaData               interface{} // meta_data (json)
	Automatic              interface{} // 0-no,1-yes
	FailureReason          interface{} //
	BillingReason          interface{} //
	Link                   interface{} //
	PaymentData            interface{} // payment create data (json)
	UniqueId               interface{} // unique id
	BalanceStart           interface{} // balance_start, utc time
	BalanceEnd             interface{} // balance_end, utc time
	InvoiceData            interface{} //
	GatewayPaymentMethod   interface{} //
	GasPayer               interface{} // who pay the gas, merchant|user
}
