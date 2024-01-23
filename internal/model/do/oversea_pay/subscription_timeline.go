// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionTimeline is the golang structure of table subscription_timeline for DAO operations like Where/Data.
type SubscriptionTimeline struct {
	g.Meta          `orm:"table:subscription_timeline, do:true"`
	Id              interface{} //
	MerchantId      interface{} // 商户Id
	UserId          interface{} // userId
	SubscriptionId  interface{} // 订阅id（内部编号）
	PeriodStart     interface{} // period_start，发票项目被添加到此发票的使用期限开始。，并非发票对应 sub 的周期
	PeriodEnd       interface{} // period_end
	PeriodStartTime *gtime.Time //
	PeriodEndTime   *gtime.Time //
	InvoiceId       interface{} // 发票ID（内部编号）
	UniqueId        interface{} // 唯一键，stripe invoice 以同步为主，其他通道 invoice 实现方案不确定，使用自定义唯一键
	Currency        interface{} // 货币
	PlanId          interface{} // 计划ID
	Quantity        interface{} // quantity
	AddonData       interface{} // plan addon json data
	ChannelId       interface{} // 支付渠道Id
	GmtCreate       *gtime.Time // 创建时间
	GmtModify       *gtime.Time // 修改时间
	IsDeleted       interface{} //
	PaymentId       interface{} // PaymentId
}
