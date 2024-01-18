// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Refund is the golang structure of table refund for DAO operations like Where/Data.
type Refund struct {
	g.Meta               `orm:"table:refund, do:true"`
	Id                   interface{} // 主键ID
	CompanyId            interface{} // 公司Id
	MerchantId           interface{} // 商户ID
	UserId               interface{} // user_id
	OpenApiId            interface{} // 使用的开放平台配置Id
	ChannelId            interface{} // 退款渠道Id
	BizType              interface{} // 业务类型。同pay.biz_type
	BizId                interface{} // 业务ID。同pay.biz_id
	CountryCode          interface{} // 国家代码，指定发起交易的国家的两个字母的ISO 3166国家代码。目前支持SG、MY、PH、ID和TH
	Currency             interface{} // 货币，“SGD” “MYR” “PHP” “IDR” “THB” 与付款金额关联的货币。指定三个字母的ISO 4217货币代码
	PaymentId            interface{} // 支付单号(内部生成，支付单号）
	RefundId             interface{} // 退款单号。可以唯一代表一笔退款（内部生成，退款单号）
	RefundFee            interface{} // 退款金额。单位：分
	RefundComment        interface{} // 退款备注
	Status               interface{} // 退款状态。10-退款中，20-退款成功，30-退款失败
	RefundTime           *gtime.Time // 退款成功时间
	GmtCreate            *gtime.Time // 创建时间
	GmtModify            *gtime.Time // 更新时间
	ChannelRefundId      interface{} // 外部退款单号
	AppId                interface{} // 退款使用的APPID
	RefundCommentExplain interface{} // 退款备注说明
	ReturnUrl            interface{} // 退款成功回调Url
	AdditionalData       interface{} //
}
