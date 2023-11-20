// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OverseaRefund is the golang structure of table oversea_refund for DAO operations like Where/Data.
type OverseaRefund struct {
	g.Meta               `orm:"table:oversea_refund, do:true"`
	Id                   interface{} // 主键ID
	CompanyId            interface{} // 公司Id
	MerchantId           interface{} // 商户ID
	BizType              interface{} // 业务类型。同pay.biz_type
	BizId                interface{} // 业务ID。同pay.biz_id
	CountryCode          interface{} // 国家代码，指定发起交易的国家的两个字母的ISO 3166国家代码。目前支持SG、MY、PH、ID和TH
	Currency             interface{} // 货币，“SGD” “MYR” “PHP” “IDR” “THB” 与付款金额关联的货币。指定三个字母的ISO 4217货币代码
	OutTradeNo           interface{} // 支付单号(内部生成，支付单号）
	OutRefundNo          interface{} // 退款单号。可以唯一代表一笔退款（内部生成，退款单号）
	RefundFee            interface{} // 退款金额。单位：分
	RefundComment        interface{} // 退款备注
	RefundStatus         interface{} // 退款状态。10-退款中，20-退款成功，30-退款失败
	RefundTime           *gtime.Time // 退款成功时间
	GmtCreate            *gtime.Time // 创建时间
	GmtModify            *gtime.Time // 更新时间
	ChannelRefundNo      interface{} // 外部退款单号
	AppId                interface{} // 退款使用的APPID
	RefundCommentExplain interface{} // 退款备注说明
	NotifyUrl            interface{} // 退款成功回调Url
	OpenApiId            interface{} // 使用的开放平台配置Id
	ChannelId            interface{} // 退款渠道Id
	ServiceRate          interface{} // 服务费比例，万分位，百分比[0，10000)，精度为0.01%，如3即为0.03%
	AdditionalData       interface{} //
}
