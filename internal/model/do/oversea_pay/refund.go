// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Refund is the golang structure of table refund for DAO operations like Where/Data.
type Refund struct {
	g.Meta               `orm:"table:refund, do:true"`
	Id                   interface{} // id
	CompanyId            interface{} // company id
	MerchantId           interface{} // merchant id
	UserId               interface{} // user_id
	OpenApiId            interface{} // open api id
	GatewayId            interface{} // gateway_id
	BizType              interface{} // biz type, copy from payment.biz_type
	ExternalRefundId     interface{} // external_refund_id
	CountryCode          interface{} // country code
	Currency             interface{} // currency
	PaymentId            interface{} // relative payment id
	RefundId             interface{} // refund id (system generate)
	RefundAmount         interface{} // refund amount, cent
	RefundComment        interface{} // refund comment
	Status               interface{} // status。10-pending，20-success，30-failure, 40-cancel
	RefundTime           interface{} // refund success time
	GmtCreate            *gtime.Time // create time
	GmtModify            *gtime.Time // update time
	GatewayRefundId      interface{} // gateway refund id
	AppId                interface{} // app id
	RefundCommentExplain interface{} // refund comment
	ReturnUrl            interface{} // return url after refund success
	MetaData             interface{} // meta_data(json)
	UniqueId             interface{} // unique id
	SubscriptionId       interface{} // subscription id
	CreateTime           interface{} // create utc time
}
