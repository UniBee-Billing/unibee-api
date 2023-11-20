// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OverseaRefund is the golang structure for table oversea_refund.
type OverseaRefund struct {
	Id                   int64       `json:"id"                   ` // 主键ID
	CompanyId            int64       `json:"companyId"            ` // 公司Id
	MerchantId           int64       `json:"merchantId"           ` // 商户ID
	BizType              int         `json:"bizType"              ` // 业务类型。同pay.biz_type
	BizId                string      `json:"bizId"                ` // 业务ID。同pay.biz_id
	CountryCode          string      `json:"countryCode"          ` // 国家代码，指定发起交易的国家的两个字母的ISO 3166国家代码。目前支持SG、MY、PH、ID和TH
	Currency             string      `json:"currency"             ` // 货币，“SGD” “MYR” “PHP” “IDR” “THB” 与付款金额关联的货币。指定三个字母的ISO 4217货币代码
	OutTradeNo           string      `json:"outTradeNo"           ` // 支付单号(内部生成，支付单号）
	OutRefundNo          string      `json:"outRefundNo"          ` // 退款单号。可以唯一代表一笔退款（内部生成，退款单号）
	RefundFee            int64       `json:"refundFee"            ` // 退款金额。单位：分
	RefundComment        string      `json:"refundComment"        ` // 退款备注
	RefundStatus         int         `json:"refundStatus"         ` // 退款状态。10-退款中，20-退款成功，30-退款失败
	RefundTime           *gtime.Time `json:"refundTime"           ` // 退款成功时间
	GmtCreate            *gtime.Time `json:"gmtCreate"            ` // 创建时间
	GmtModify            *gtime.Time `json:"gmtModify"            ` // 更新时间
	ChannelRefundNo      string      `json:"channelRefundNo"      ` // 外部退款单号
	AppId                string      `json:"appId"                ` // 退款使用的APPID
	RefundCommentExplain string      `json:"refundCommentExplain" ` // 退款备注说明
	NotifyUrl            string      `json:"notifyUrl"            ` // 退款成功回调Url
	OpenApiId            int64       `json:"openApiId"            ` // 使用的开放平台配置Id
	ChannelId            int64       `json:"channelId"            ` // 退款渠道Id
	ServiceRate          int64       `json:"serviceRate"          ` // 服务费比例，万分位，百分比[0，10000)，精度为0.01%，如3即为0.03%
	AdditionalData       string      `json:"additionalData"       ` //
}
