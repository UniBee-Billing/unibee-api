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
	Code               interface{} //
	Status             interface{} // status, 1-editable, 2-active, 3-deactive, 4-expire
	BillingType        interface{} // billing_type, 1-one-time, 2-recurring
	DiscountType       interface{} // discount_type, 1-percentage, 2-fixed_amount
	DiscountAmount     interface{} // amount of discount, available when discount_type is fixed_amount
	DiscountPercentage interface{} // percentage of discount, 100=1%, available when discount_type is percentage
	Currency           interface{} // currency of discount, available when discount_type is fixed_amount
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
	PlanIds            interface{} // Ids of plan which discount code can effect, default effect all plans if not set
	Quantity           interface{} // quantity of code
	Advance            interface{} // AdvanceConfig,  0-false,1-true, will enable all advance config if set 1
	UserLimit          interface{} // AdvanceConfig, The limit of every customer can apply, the recurring apply not involved, 0-unlimited"
	UserScope          interface{} // AdvanceConfig, Apply user scope,0-for all, 1-for only new user, 2-for only renewals, renewals is upgrade&downgrade&renew
	UpgradeOnly        interface{} // AdvanceConfig, 0-false,1-true, will forbid for all except upgrade action if set 1
	UpgradeLongerOnly  interface{} // AdvanceConfig, 0-false,1-true, will forbid for all except upgrade to longer plan if set 1
	PlanApplyType      interface{} // plan apply type, 0-apply for all, 1-apply for plans specified, 2-exclude for plans specified
}
