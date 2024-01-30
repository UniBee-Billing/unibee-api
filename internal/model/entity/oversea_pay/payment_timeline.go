// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentTimeline is the golang structure for table payment_timeline.
type PaymentTimeline struct {
	Id             uint64      `json:"id"             description:""`                                                       //
	MerchantId     int64       `json:"merchantId"     description:"商户Id"`                                                   // 商户Id
	UserId         int64       `json:"userId"         description:"userId"`                                                 // userId
	SubscriptionId string      `json:"subscriptionId" description:"订阅id（内部编号）"`                                             // 订阅id（内部编号）
	InvoiceId      string      `json:"invoiceId"      description:"发票ID（内部编号）"`                                             // 发票ID（内部编号）
	UniqueId       string      `json:"uniqueId"       description:"唯一键，stripe invoice 以同步为主，其他通道 invoice 实现方案不确定，使用自定义唯一键"` // 唯一键，stripe invoice 以同步为主，其他通道 invoice 实现方案不确定，使用自定义唯一键
	Currency       string      `json:"currency"       description:"货币"`                                                     // 货币
	TotalAmount    int64       `json:"totalAmount"    description:"金额,单位：分"`                                                // 金额,单位：分
	ChannelId      int64       `json:"channelId"      description:"支付渠道Id"`                                                 // 支付渠道Id
	GmtCreate      *gtime.Time `json:"gmtCreate"      description:"创建时间"`                                                   // 创建时间
	GmtModify      *gtime.Time `json:"gmtModify"      description:"修改时间"`                                                   // 修改时间
	IsDeleted      int         `json:"isDeleted"      description:"0-UnDeleted，1-Deleted"`                                  // 0-UnDeleted，1-Deleted
	PaymentId      string      `json:"paymentId"      description:"PaymentId"`                                              // PaymentId
	Status         int         `json:"status"         description:"0-pending, 1-success, 2-failure"`                        // 0-pending, 1-success, 2-failure
	TimelineType   int         `json:"timelineType"   description:"0-pay, 1-refund"`                                        // 0-pay, 1-refund
}
