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
	TaxPercentage             interface{} // taxPercentage 1000 = 10%
	TaxInclusive              interface{} // deperated
	Type                      interface{} // type，1-main plan，2-recurring addon plan 3-onetime addon plan
	Status                    interface{} // status，1-editing，2-active，3-inactive，4-expired
	IsDeleted                 interface{} // 0-UnDeleted，1-Deleted
	BindingAddonIds           interface{} // binded recurring addon planIds，split with ,
	BindingOnetimeAddonIds    interface{} // binded onetime addon planIds，split with ,
	PublishStatus             interface{} // 1-UnPublish,2-Publish, Use For Display Plan At UserPortal
	CreateTime                interface{} // create utc time
	ExtraMetricData           interface{} //
	MetaData                  interface{} // meta_data(json)
	GasPayer                  interface{} // who pay the gas, merchant|user
	TrialAmount               interface{} // amount of trial, 0 for free
	TrialDurationTime         interface{} // duration of trial
	TrialDemand               interface{} //
	CancelAtTrialEnd          interface{} // whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription
}
