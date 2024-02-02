// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionPlan is the golang structure of table subscription_plan for DAO operations like Where/Data.
type SubscriptionPlan struct {
	g.Meta                    `orm:"table:subscription_plan, do:true"`
	Id                        interface{} //
	GmtCreate                 *gtime.Time // create time
	GmtModify                 *gtime.Time // update time
	CompanyId                 interface{} // company id
	MerchantId                interface{} // merchant id
	PlanName                  interface{} // PlanName
	Amount                    interface{} // amount, cent, without tax
	Currency                  interface{} // currency
	IntervalUnit              interface{} // period unit,day|month|year|week
	IntervalCount             interface{} // period unit count
	Description               interface{} // description
	ImageUrl                  interface{} // image_url
	HomeUrl                   interface{} // home_url
	ChannelProductName        interface{} // channel product name
	ChannelProductDescription interface{} // channel product description
	TaxScale                  interface{} // tax scale 1000 = 10%
	TaxInclusive              interface{} // deperated
	Type                      interface{} // type，1-main plan，2-addon plan
	Status                    interface{} // status，1-editing，2-active，3-inactive，4-expired
	IsDeleted                 interface{} // 0-UnDeleted，1-Deleted
	BindingAddonIds           interface{} // binded addon planIds，split with ,
	PublishStatus             interface{} // 1-UnPublish,2-Publish,用于控制是否在 UserPortal 端展示
}
