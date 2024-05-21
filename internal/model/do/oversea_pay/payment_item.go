// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentItem is the golang structure of table payment_item for DAO operations like Where/Data.
type PaymentItem struct {
	g.Meta         `orm:"table:payment_item, do:true"`
	Id             interface{} //
	BizType        interface{} // biz_type 1-onetime payment, 3-subscription
	MerchantId     interface{} // merchant id
	UserId         interface{} // userId
	SubscriptionId interface{} // subscription id
	InvoiceId      interface{} // invoice id
	UniqueId       interface{} // unique id
	Currency       interface{} // currency
	Amount         interface{} // amount
	UnitAmount     interface{} // unit_amount
	Quantity       interface{} // quantity
	GatewayId      interface{} // gateway id
	GmtCreate      *gtime.Time // create time
	GmtModify      *gtime.Time // update time
	IsDeleted      interface{} // 0-UnDeleted，1-Deleted
	PaymentId      interface{} // PaymentId
	Status         interface{} // 0-pending, 1-success, 2-failure
	CreateTime     interface{} // create utc time
	Description    interface{} // description
	Name           interface{} // name
}
