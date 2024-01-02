// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionPlan is the golang structure for table subscription_plan.
type SubscriptionPlan struct {
	Id                        uint64      `json:"id"                        ` //
	GmtCreate                 *gtime.Time `json:"gmtCreate"                 ` // 创建时间
	GmtModify                 *gtime.Time `json:"gmtModify"                 ` // 修改时间
	CompanyId                 int64       `json:"companyId"                 ` // 公司ID
	MerchantId                int64       `json:"merchantId"                ` // 商户Id
	PlanName                  string      `json:"planName"                  ` // 计划名称
	Amount                    int64       `json:"amount"                    ` // 金额,单位：分
	Currency                  string      `json:"currency"                  ` // 货币
	IntervalUnit              string      `json:"intervalUnit"              ` // 周期-全小写存放,day|month|year|week
	IntervalCount             int         `json:"intervalCount"             ` // 订阅计费之间的间隔数。例如，每 3 个月interval=month计费一次interval_count=3。允许的最长间隔为一年（1 年、12 个月或 52 周）
	Description               string      `json:"description"               ` //
	IsDeleted                 int         `json:"isDeleted"                 ` //
	ImageUrl                  string      `json:"imageUrl"                  ` // image_url
	HomeUrl                   string      `json:"homeUrl"                   ` // home_url
	ChannelProductName        string      `json:"channelProductName"        ` // 支付渠道product_name
	ChannelProductDescription string      `json:"channelProductDescription" ` // 支付渠道product_description
	TaxPercentage             int         `json:"taxPercentage"             ` // 税费比例： 1 =1%
	TaxInclusive              int         `json:"taxInclusive"              ` // 税费是否包含，1-包含，0-不包含
	Type                      int         `json:"type"                      ` // 类型，1-main plan，2-addon plan
	Status                    int         `json:"status"                    ` // 状态，1-编辑中，2-活跃，3-非活跃，4-过期
	BindingAddonIds           string      `json:"bindingAddonIds"           ` // 绑定的 Addon PlanIds，以逗号隔开
}
