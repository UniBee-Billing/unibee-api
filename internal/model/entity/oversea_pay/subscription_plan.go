// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionPlan is the golang structure for table subscription_plan.
type SubscriptionPlan struct {
	Id           uint64      `json:"id"           ` //
	GmtCreate    *gtime.Time `json:"gmtCreate"    ` // 创建时间
	GmtModify    *gtime.Time `json:"gmtModify"    ` // 修改时间
	CompanyId    int64       `json:"companyId"    ` // 公司ID
	MerchantId   int64       `json:"merchantId"   ` // 商户Id
	PlanName     string      `json:"planName"     ` // 计划名称
	Amount       string      `json:"amount"       ` // 金额,单位：分
	Currency     string      `json:"currency"     ` // 货币
	IntervalUnit string      `json:"intervalUnit" ` // 周期,day|month|year|week
	Description  string      `json:"description"  ` //
	IsDeleted    int         `json:"isDeleted"    ` //
	ImageUrl     string      `json:"imageUrl"     ` // image_url
	HomeUrl      string      `json:"homeUrl"      ` // home_url
}
