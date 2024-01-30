// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentTimeline is the golang structure of table payment_timeline for DAO operations like Where/Data.
type PaymentTimeline struct {
	g.Meta         `orm:"table:payment_timeline, do:true"`
	Id             interface{} //
	MerchantId     interface{} // 商户Id
	UserId         interface{} // userId
	SubscriptionId interface{} // 订阅id（内部编号）
	InvoiceId      interface{} // 发票ID（内部编号）
	UniqueId       interface{} // 唯一键，stripe invoice 以同步为主，其他通道 invoice 实现方案不确定，使用自定义唯一键
	Currency       interface{} // 货币
	TotalAmount    interface{} // 金额,单位：分
	ChannelId      interface{} // 支付渠道Id
	GmtCreate      *gtime.Time // 创建时间
	GmtModify      *gtime.Time // 修改时间
	IsDeleted      interface{} // 0-UnDeleted，1-Deleted
	PaymentId      interface{} // PaymentId
	Status         interface{} // 0-pending, 1-success, 2-failure
	TimelineType   interface{} // 0-pay, 1-refund
}
