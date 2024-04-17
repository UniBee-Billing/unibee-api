// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantDiscountCode is the golang structure of table merchant_discount_code for DAO operations like Where/Data.
type MerchantDiscountCode struct {
	g.Meta             `orm:"table:merchant_discount_code, do:true"`
	Id                 interface{} // ID
	MerchantId         interface{} // merchantId
	Name               interface{} // name
	Code               interface{} // code
	Status             interface{} // status, 1-editable, 2-active, 3-deactive, 4-expire
	BillingType        interface{} // billing_type, 1-one-time, 2-recurring
	DiscountType       interface{} // discount_type, 1-percentage, 2-fixed_amount
	DiscountAmount     interface{} // amount of discount, available when discount_type is fixed_amount
	DiscountPercentage interface{} // percentage of discount, 100=1%, available when discount_type is percentage
	Currency           interface{} // currency of discount, available when discount_type is fixed_amount
	UserLimit          interface{} // the limit of every user apply, 0-unlimited
	SubscriptionLimit  interface{} // the limit of every subscription apply, 0-unlimited
	StartTime          interface{} // start of discount available utc time
	EndTime            interface{} // end of discount available utc time, 0-invalid
	GmtCreate          *gtime.Time // create time
	GmtModify          *gtime.Time // update time
	IsDeleted          interface{} // 0-UnDeleted，1-Deleted
	CreateTime         interface{} // create utc time
	CycleLimit         interface{} // the count limitation of subscription cycle , 0-no limit
	MetaData           interface{} // meta_data(json)
	Type               interface{} // type, 1-external discount code
}
