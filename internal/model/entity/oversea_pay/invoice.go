// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Invoice is the golang structure for table invoice.
type Invoice struct {
	Id                             uint64      `json:"id"                             description:""`                                                                      //
	MerchantId                     int64       `json:"merchantId"                     description:"商户Id"`                                                                  // 商户Id
	UserId                         int64       `json:"userId"                         description:"userId"`                                                                // userId
	SubscriptionId                 string      `json:"subscriptionId"                 description:"订阅id（内部编号）"`                                                            // 订阅id（内部编号）
	InvoiceId                      string      `json:"invoiceId"                      description:"发票ID（内部编号）"`                                                            // 发票ID（内部编号）
	InvoiceName                    string      `json:"invoiceName"                    description:"发票名称"`                                                                  // 发票名称
	UniqueId                       string      `json:"uniqueId"                       description:"唯一键，stripe invoice 以同步为主，其他通道 invoice 实现方案不确定，使用自定义唯一键"`                // 唯一键，stripe invoice 以同步为主，其他通道 invoice 实现方案不确定，使用自定义唯一键
	GmtCreate                      *gtime.Time `json:"gmtCreate"                      description:"创建时间"`                                                                  // 创建时间
	TotalAmount                    int64       `json:"totalAmount"                    description:"金额,单位：分"`                                                               // 金额,单位：分
	TaxAmount                      int64       `json:"taxAmount"                      description:"Tax金额,单位：分"`                                                            // Tax金额,单位：分
	SubscriptionAmount             int64       `json:"subscriptionAmount"             description:"Sub金额,单位：分"`                                                            // Sub金额,单位：分
	Currency                       string      `json:"currency"                       description:"货币"`                                                                    // 货币
	Lines                          string      `json:"lines"                          description:"lines json data"`                                                       // lines json data
	ChannelId                      int64       `json:"channelId"                      description:"支付渠道Id"`                                                                // 支付渠道Id
	Status                         int         `json:"status"                         description:"订阅单状态，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled"` // 订阅单状态，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
	SendStatus                     int         `json:"sendStatus"                     description:"邮件发送状态，0-No | 1- YES"`                                                  // 邮件发送状态，0-No | 1- YES
	SendEmail                      string      `json:"sendEmail"                      description:"email 发送地址，取自 UserAccount 表 email"`                                     // email 发送地址，取自 UserAccount 表 email
	SendPdf                        string      `json:"sendPdf"                        description:"pdf 文件地址"`                                                              // pdf 文件地址
	GmtModify                      *gtime.Time `json:"gmtModify"                      description:"修改时间"`                                                                  // 修改时间
	IsDeleted                      int         `json:"isDeleted"                      description:""`                                                                      //
	Link                           string      `json:"link"                           description:"invoice 链接（可用于支付）"`                                                     // invoice 链接（可用于支付）
	ChannelStatus                  string      `json:"channelStatus"                  description:"渠道最新状态，Stripe：https://stripe.com/docs/api/invoices/object"`             // 渠道最新状态，Stripe：https://stripe.com/docs/api/invoices/object
	ChannelPaymentId               string      `json:"channelPaymentId"               description:"关联渠道 PaymentId"`                                                        // 关联渠道 PaymentId
	ChannelUserId                  string      `json:"channelUserId"                  description:"渠道用户 Id"`                                                               // 渠道用户 Id
	ChannelInvoicePdf              string      `json:"channelInvoicePdf"              description:"关联渠道发票 pdf"`                                                            // 关联渠道发票 pdf
	TaxPercentage                  int64       `json:"taxPercentage"                  description:"Tax税率，万分位，1000 表示 10%"`                                                 // Tax税率，万分位，1000 表示 10%
	SendNote                       string      `json:"sendNote"                       description:"send_note"`                                                             // send_note
	SendTerms                      string      `json:"sendTerms"                      description:"send_terms"`                                                            // send_terms
	TotalAmountExcludingTax        int64       `json:"totalAmountExcludingTax"        description:"金额(不含税）,单位：分"`                                                          // 金额(不含税）,单位：分
	SubscriptionAmountExcludingTax int64       `json:"subscriptionAmountExcludingTax" description:"Sub金额(不含税）,单位：分"`                                                       // Sub金额(不含税）,单位：分
	ChannelInvoiceId               string      `json:"channelInvoiceId"               description:"关联渠道发票 Id"`                                                             // 关联渠道发票 Id
	PeriodStart                    int64       `json:"periodStart"                    description:"period_start，发票项目被添加到此发票的使用期限开始。，并非发票对应 sub 的周期"`                       // period_start，发票项目被添加到此发票的使用期限开始。，并非发票对应 sub 的周期
	PeriodEnd                      int64       `json:"periodEnd"                      description:"period_end"`                                                            // period_end
	PeriodStartTime                *gtime.Time `json:"periodStartTime"                description:""`                                                                      //
	PeriodEndTime                  *gtime.Time `json:"periodEndTime"                  description:""`                                                                      //
	PaymentId                      string      `json:"paymentId"                      description:"PaymentId"`                                                             // PaymentId
	RefundId                       string      `json:"refundId"                       description:"refundId"`                                                              // refundId
	Data                           string      `json:"data"                           description:"渠道额外参数，JSON格式"`                                                         // 渠道额外参数，JSON格式
}
