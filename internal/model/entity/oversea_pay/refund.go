// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Refund is the golang structure for table refund.
type Refund struct {
	Id                   int64       `json:"id"                   description:"主键ID"`                                                            // 主键ID
	CompanyId            int64       `json:"companyId"            description:"公司Id"`                                                            // 公司Id
	MerchantId           int64       `json:"merchantId"           description:"商户ID"`                                                            // 商户ID
	UserId               int64       `json:"userId"               description:"user_id"`                                                         // user_id
	OpenApiId            int64       `json:"openApiId"            description:"使用的开放平台配置Id"`                                                     // 使用的开放平台配置Id
	ChannelId            int64       `json:"channelId"            description:"退款渠道Id"`                                                          // 退款渠道Id
	BizType              int         `json:"bizType"              description:"业务类型。同pay.biz_type"`                                              // 业务类型。同pay.biz_type
	BizId                string      `json:"bizId"                description:"业务ID。同pay.biz_id"`                                                // 业务ID。同pay.biz_id
	CountryCode          string      `json:"countryCode"          description:"国家代码，指定发起交易的国家的两个字母的ISO 3166国家代码。目前支持SG、MY、PH、ID和TH"`             // 国家代码，指定发起交易的国家的两个字母的ISO 3166国家代码。目前支持SG、MY、PH、ID和TH
	Currency             string      `json:"currency"             description:"货币，“SGD” “MYR” “PHP” “IDR” “THB” 与付款金额关联的货币。指定三个字母的ISO 4217货币代码"` // 货币，“SGD” “MYR” “PHP” “IDR” “THB” 与付款金额关联的货币。指定三个字母的ISO 4217货币代码
	PaymentId            string      `json:"paymentId"            description:"支付单号(内部生成，支付单号）"`                                                 // 支付单号(内部生成，支付单号）
	RefundId             string      `json:"refundId"             description:"退款单号。可以唯一代表一笔退款（内部生成，退款单号）"`                                      // 退款单号。可以唯一代表一笔退款（内部生成，退款单号）
	RefundAmount         int64       `json:"refundAmount"         description:"退款金额。单位：分"`                                                       // 退款金额。单位：分
	RefundComment        string      `json:"refundComment"        description:"退款备注"`                                                            // 退款备注
	Status               int         `json:"status"               description:"退款状态。10-退款中，20-退款成功，30-退款失败"`                                     // 退款状态。10-退款中，20-退款成功，30-退款失败
	RefundTime           *gtime.Time `json:"refundTime"           description:"退款成功时间"`                                                          // 退款成功时间
	GmtCreate            *gtime.Time `json:"gmtCreate"            description:"创建时间"`                                                            // 创建时间
	GmtModify            *gtime.Time `json:"gmtModify"            description:"更新时间"`                                                            // 更新时间
	ChannelRefundId      string      `json:"channelRefundId"      description:"外部退款单号"`                                                          // 外部退款单号
	AppId                string      `json:"appId"                description:"退款使用的APPID"`                                                      // 退款使用的APPID
	RefundCommentExplain string      `json:"refundCommentExplain" description:"退款备注说明"`                                                          // 退款备注说明
	ReturnUrl            string      `json:"returnUrl"            description:"退款成功回调Url"`                                                       // 退款成功回调Url
	AdditionalData       string      `json:"additionalData"       description:""`                                                                //
	UniqueId             string      `json:"uniqueId"             description:"唯一键，以同步为逻辑加入使用自定义唯一键"`                                            // 唯一键，以同步为逻辑加入使用自定义唯一键
	SubscriptionId       string      `json:"subscriptionId"       description:"订阅id（内部编号）"`                                                      // 订阅id（内部编号）
}
