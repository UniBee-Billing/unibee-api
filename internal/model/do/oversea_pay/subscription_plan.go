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
	g.Meta       `orm:"table:subscription_plan, do:true"`
	Id           interface{} //
	GmtCreate    *gtime.Time // 创建时间
	GmtModify    *gtime.Time // 修改时间
	CompanyId    interface{} // 公司ID
	MerchantId   interface{} // 商户Id
	PlanName     interface{} // 计划名称
	Amount       interface{} // 金额,单位：分
	Currency     interface{} // 货币
	IntervalUnit interface{} // 周期,day|month|year|week
	Description  interface{} //
	IsDeleted    interface{} //
}
