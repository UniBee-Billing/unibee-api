// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Plan is the golang structure of table plan for DAO operations like Where/Data.
type Plan struct {
	g.Meta                    `orm:"table:plan, do:true"`
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
	GatewayProductName        interface{} // gateway product name
	GatewayProductDescription interface{} // gateway product description
	TaxScale                  interface{} // tax scale 1000 = 10%
	TaxInclusive              interface{} // deperated
	Type                      interface{} // type，1-main plan，2-addon plan
	Status                    interface{} // status，1-editing，2-active，3-inactive，4-expired
	IsDeleted                 interface{} // 0-UnDeleted，1-Deleted
	BindingAddonIds           interface{} // binded addon planIds，split with ,
	PublishStatus             interface{} // 1-UnPublish,2-Publish, Use For Display Plan At UserPortal
	CreateTime                interface{} // create utc time
}