// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantUserDiscountCode is the golang structure of table merchant_user_discount_code for DAO operations like Where/Data.
type MerchantUserDiscountCode struct {
	g.Meta         `orm:"table:merchant_user_discount_code, do:true"`
	Id             interface{} // ID
	MerchantId     interface{} // merchantId
	UserId         interface{} // user_id
	Code           interface{} // code
	Status         interface{} // status, 1-normal, 2-rollback
	PlanId         interface{} // plan_id
	SubscriptionId interface{} // subscription_id
	PaymentId      interface{} // payment_id
	InvoiceId      interface{} // invoice_id
	UniqueId       interface{} // unique_id
	GmtCreate      *gtime.Time // create time
	GmtModify      *gtime.Time // update time
	IsDeleted      interface{} // 0-UnDeleted，1-Deleted
	CreateTime     interface{} // create utc time
}
