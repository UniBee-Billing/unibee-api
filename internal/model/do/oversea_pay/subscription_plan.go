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
	GmtCreate                 *gtime.Time // 创建时间
	GmtModify                 *gtime.Time // 修改时间
	CompanyId                 interface{} // 公司ID
	MerchantId                interface{} // 商户Id
	PlanName                  interface{} // 计划名称
	Amount                    interface{} // 金额,单位：分
	Currency                  interface{} // 货币
	IntervalUnit              interface{} // 周期-全小写存放,day|month|year|week
	Description               interface{} //
	IsDeleted                 interface{} //
	ImageUrl                  interface{} // image_url
	HomeUrl                   interface{} // home_url
	ChannelProductName        interface{} // 支付渠道product_name
	ChannelProductDescription interface{} // 支付渠道product_description
	TaxPercentage             interface{} // 税费比例： 1 =1%
	TaxInclusive              interface{} // 税费是否包含，1-包含，0-不包含
}
