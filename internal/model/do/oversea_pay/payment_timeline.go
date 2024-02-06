// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentTimeline is the golang structure of table payment_timeline for DAO operations like Where/Data.
type PaymentTimeline struct {
	g.Meta         `orm:"table:payment_timeline, do:true"`
	Id             interface{} //
	MerchantId     interface{} // merchant id
	UserId         interface{} // userId
	SubscriptionId interface{} // subscription id
	InvoiceId      interface{} // invoice id
	UniqueId       interface{} // unique id
	Currency       interface{} // currency
	TotalAmount    interface{} // total amount
	GatewayId      interface{} // gateway id
	GmtCreate      *gtime.Time // create time
	GmtModify      *gtime.Time // update time
	IsDeleted      interface{} // 0-UnDeleted，1-Deleted
	PaymentId      interface{} // PaymentId
	Status         interface{} // 0-pending, 1-success, 2-failure
	TimelineType   interface{} // 0-pay, 1-refund
	CreateAt       interface{} // create utc time
}
