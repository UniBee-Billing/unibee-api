// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionOnetimeAddon is the golang structure of table subscription_onetime_addon for DAO operations like Where/Data.
type SubscriptionOnetimeAddon struct {
	g.Meta         `orm:"table:subscription_onetime_addon, do:true"`
	Id             interface{} // id
	GmtCreate      *gtime.Time // create_time
	GmtModify      *gtime.Time // modify_time
	SubscriptionId interface{} // subscription_id
	AddonId        interface{} // onetime addonId
	Quantity       interface{} // quantity
	Status         interface{} // status, 1-create, 2-paid, 3-cancel, 4-expired
	IsDeleted      interface{} // 0-UnDeleted，1-Deleted
	CreateTime     interface{} // create utc time
	PaymentId      interface{} // paymentId
	MetaData       interface{} // meta_data(json)
	UserId         interface{} // userId
}
