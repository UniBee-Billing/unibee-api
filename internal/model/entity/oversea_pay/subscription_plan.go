// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionPlan is the golang structure for table subscription_plan.
type SubscriptionPlan struct {
	Id                        uint64      `json:"id"                        description:""`                                                                                   //
	GmtCreate                 *gtime.Time `json:"gmtCreate"                 description:"创建时间"`                                                                               // 创建时间
	GmtModify                 *gtime.Time `json:"gmtModify"                 description:"修改时间"`                                                                               // 修改时间
	CompanyId                 int64       `json:"companyId"                 description:"公司ID"`                                                                               // 公司ID
	MerchantId                int64       `json:"merchantId"                description:"商户Id"`                                                                               // 商户Id
	PlanName                  string      `json:"planName"                  description:"计划名称"`                                                                               // 计划名称
	Amount                    int64       `json:"amount"                    description:"金额,单位：分"`                                                                            // 金额,单位：分
	Currency                  string      `json:"currency"                  description:"货币"`                                                                                 // 货币
	IntervalUnit              string      `json:"intervalUnit"              description:"周期-全小写存放,day|month|year|week"`                                                       // 周期-全小写存放,day|month|year|week
	IntervalCount             int         `json:"intervalCount"             description:"订阅计费之间的间隔数。例如，每 3 个月interval=month计费一次interval_count=3。允许的最长间隔为一年（1 年、12 个月或 52 周）"` // 订阅计费之间的间隔数。例如，每 3 个月interval=month计费一次interval_count=3。允许的最长间隔为一年（1 年、12 个月或 52 周）
	Description               string      `json:"description"               description:""`                                                                                   //
	IsDeleted                 int         `json:"isDeleted"                 description:""`                                                                                   //
	ImageUrl                  string      `json:"imageUrl"                  description:"image_url"`                                                                          // image_url
	HomeUrl                   string      `json:"homeUrl"                   description:"home_url"`                                                                           // home_url
	ChannelProductName        string      `json:"channelProductName"        description:"支付渠道product_name"`                                                                   // 支付渠道product_name
	ChannelProductDescription string      `json:"channelProductDescription" description:"支付渠道product_description"`                                                            // 支付渠道product_description
	TaxPercentage             int         `json:"taxPercentage"             description:"税费比例： 1 =1%"`                                                                        // 税费比例： 1 =1%
	TaxInclusive              int         `json:"taxInclusive"              description:"税费是否包含，1-包含，0-不包含"`                                                                  // 税费是否包含，1-包含，0-不包含
	Type                      int         `json:"type"                      description:"类型，1-main plan，2-addon plan"`                                                        // 类型，1-main plan，2-addon plan
	Status                    int         `json:"status"                    description:"状态，1-编辑中，2-活跃，3-非活跃，4-过期"`                                                           // 状态，1-编辑中，2-活跃，3-非活跃，4-过期
	BindingAddonIds           string      `json:"bindingAddonIds"           description:"绑定的 Addon PlanIds，以逗号隔开"`                                                            // 绑定的 Addon PlanIds，以逗号隔开
	PublishStatus             int         `json:"publishStatus"             description:"1-UnPublish,2-Publish,用于控制是否在 UserPortal 端展示"`                                       // 1-UnPublish,2-Publish,用于控制是否在 UserPortal 端展示
}
