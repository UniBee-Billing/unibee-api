// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Invoice is the golang structure of table invoice for DAO operations like Where/Data.
type Invoice struct {
	g.Meta                         `orm:"table:invoice, do:true"`
	Id                             interface{} //
	MerchantId                     interface{} // 商户Id
	UserId                         interface{} // userId
	SubscriptionId                 interface{} // 订阅id（内部编号）
	InvoiceId                      interface{} // 发票ID（内部编号）
	InvoiceName                    interface{} // 发票名称
	UniqueId                       interface{} // 唯一键，stripe invoice 以同步为主，其他通道 invoice 实现方案不确定，使用自定义唯一键
	GmtCreate                      *gtime.Time // 创建时间
	TotalAmount                    interface{} // 金额,单位：分
	TaxAmount                      interface{} // Tax金额,单位：分
	SubscriptionAmount             interface{} // Sub金额,单位：分
	Currency                       interface{} // 货币
	Lines                          interface{} // lines json data
	ChannelId                      interface{} // 支付渠道Id
	Status                         interface{} // 订阅单状态，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
	SendStatus                     interface{} // 邮件发送状态，0-No | 1- YES
	SendEmail                      interface{} // email 发送地址，取自 UserAccount 表 email
	SendPdf                        interface{} // pdf 文件地址
	GmtModify                      *gtime.Time // 修改时间
	IsDeleted                      interface{} //
	Link                           interface{} // invoice 链接（可用于支付）
	ChannelStatus                  interface{} // 渠道最新状态，Stripe：https://stripe.com/docs/api/invoices/object
	ChannelInvoiceId               interface{} // 关联渠道发票 Id
	ChannelPaymentId               interface{} // 关联渠道 PaymentId
	ChannelUserId                  interface{} // 渠道用户 Id-废弃
	ChannelInvoicePdf              interface{} // 关联渠道发票 pdf
	TaxScale                       interface{} // Tax税率，万分位，1000 表示 10%
	SendNote                       interface{} // send_note
	SendTerms                      interface{} // send_terms
	TotalAmountExcludingTax        interface{} // 金额(不含税）,单位：分
	SubscriptionAmountExcludingTax interface{} // Sub金额(不含税）,单位：分
	PeriodStart                    interface{} // period_start，发票项目被添加到此发票的使用期限开始。，并非发票对应 sub 的周期
	PeriodEnd                      interface{} // period_end
	PeriodStartTime                *gtime.Time //
	PeriodEndTime                  *gtime.Time //
	PaymentId                      interface{} // PaymentId
	RefundId                       interface{} // refundId
	Data                           interface{} // 渠道额外参数，JSON格式
	BizType                        interface{} // 业务类型。1-single payment, 3-subscription
}
