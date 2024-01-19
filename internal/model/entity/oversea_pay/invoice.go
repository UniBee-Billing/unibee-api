// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Invoice is the golang structure for table invoice.
type Invoice struct {
	Id                             uint64      `json:"id"                             ` //
	MerchantId                     int64       `json:"merchantId"                     ` // 商户Id
	UserId                         int64       `json:"userId"                         ` // userId
	SubscriptionId                 string      `json:"subscriptionId"                 ` // 订阅id（内部编号）
	InvoiceId                      string      `json:"invoiceId"                      ` // 发票ID（内部编号）
	ChannelInvoiceId               string      `json:"channelInvoiceId"               ` // 关联渠道发票 Id
	UniqueId                       string      `json:"uniqueId"                       ` // 唯一键，stripe invoice 以同步为主，其他通道 invoice 实现方案不确定，使用自定义唯一键
	GmtCreate                      *gtime.Time `json:"gmtCreate"                      ` // 创建时间
	TotalAmount                    int64       `json:"totalAmount"                    ` // 金额,单位：分
	TaxAmount                      int64       `json:"taxAmount"                      ` // Tax金额,单位：分
	SubscriptionAmount             int64       `json:"subscriptionAmount"             ` // Sub金额,单位：分
	Currency                       string      `json:"currency"                       ` // 货币
	Lines                          string      `json:"lines"                          ` // lines json data
	ChannelId                      int64       `json:"channelId"                      ` // 支付渠道Id
	Status                         int         `json:"status"                         ` // 订阅单状态，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
	SendStatus                     int         `json:"sendStatus"                     ` // 邮件发送状态，0-No | 1- YES
	SendEmail                      string      `json:"sendEmail"                      ` // email 发送地址，取自 UserAccount 表 email
	SendPdf                        string      `json:"sendPdf"                        ` // pdf 文件地址
	Data                           string      `json:"data"                           ` // 渠道额外参数，JSON格式
	GmtModify                      *gtime.Time `json:"gmtModify"                      ` // 修改时间
	IsDeleted                      int         `json:"isDeleted"                      ` //
	Link                           string      `json:"link"                           ` // invoice 链接（可用于支付）
	ChannelStatus                  string      `json:"channelStatus"                  ` // 渠道最新状态，Stripe：https://stripe.com/docs/api/invoices/object
	ChannelPaymentId               string      `json:"channelPaymentId"               ` // 关联渠道 PaymentId
	ChannelUserId                  string      `json:"channelUserId"                  ` // 渠道用户 Id
	ChannelInvoicePdf              string      `json:"channelInvoicePdf"              ` // 关联渠道发票 pdf
	TaxPercentage                  int64       `json:"taxPercentage"                  ` // Tax税率，万分位，1000 表示 10%
	SendNote                       string      `json:"sendNote"                       ` // send_note
	SendTerms                      string      `json:"sendTerms"                      ` // send_terms
	TotalAmountExcludingTax        int64       `json:"totalAmountExcludingTax"        ` // 金额(不含税）,单位：分
	SubscriptionAmountExcludingTax int64       `json:"subscriptionAmountExcludingTax" ` // Sub金额(不含税）,单位：分
	PeriodStart                    int64       `json:"periodStart"                    ` // period_start
	PeriodEnd                      int64       `json:"periodEnd"                      ` // period_end
	PaymentId                      string      `json:"paymentId"                      ` // PaymentId
	RefundId                       string      `json:"refundId"                       ` // refundId
}
