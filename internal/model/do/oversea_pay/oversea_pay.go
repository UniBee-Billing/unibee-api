// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OverseaPay is the golang structure of table oversea_pay for DAO operations like Where/Data.
type OverseaPay struct {
	g.Meta             `orm:"table:oversea_pay, do:true"`
	Id                 interface{} // 主键id
	CompanyId          interface{} // 公司Id
	MerchantId         interface{} // 商户ID
	BizType            interface{} // 业务类型。1-订单
	BizId              interface{} // 业务id-即商户订单号
	CountryCode        interface{} // 国家代码，指定发起交易的国家的两个字母的ISO 3166国家代码。目前支持SG、MY、PH、ID和TH
	Currency           interface{} // 货币，“SGD” “MYR” “PHP” “IDR” “THB” 与付款金额关联的货币。指定三个字母的ISO 4217货币代码
	MerchantOrderNo    interface{} // 内部支付编号（系统生成唯一）
	PaymentFee         interface{} // 支付金额
	RefundFee          interface{} // 总共已退款金额
	BuyerPayFee        interface{} // 买家实付金额
	ReceiptFee         interface{} // 商户捕获金额（分）
	PayStatus          interface{} // 支付状态。10-支付中，20-支付成功，30-支付取消
	AuthorizeStatus    interface{} // 用户授权状态，0-未授权，1-已授权，2-已发起捕获
	ChannelId          interface{} // 支付方式id,表oversea_pay_channel的id
	CaptureDelayHours  interface{} // 延迟Capture时间
	ChannelPayId       interface{} // 第三方支付平台支付预订单ID，支付接口返回
	ChannelTradeNo     interface{} // 外部支付渠道订单号，支付成功回调返回
	CreateTime         *gtime.Time // 支付单创建时间
	CancelTime         *gtime.Time // 支付单取消时间
	PaidTime           *gtime.Time // 付款成功时间
	InvoiceTime        *gtime.Time // 入账成功时间
	GmtCreate          *gtime.Time // 创建时间
	GmtModify          *gtime.Time // 更新时间
	InvoiceStatus      interface{} // 入账状态，未入账-0，入账中-1，完成入账-2，入账失败-3，已撤销入账-4
	InvoiceFee         interface{} // 入账总金额。单位：分，invoice_total_fee + service_fee = payment_fee - refund_fee
	ServiceRate        interface{} // 服务费比例，万分位，百分比[0，10000)，精度为0.01%，如3即为0.03%
	ServiceFee         interface{} // 服务费。单位：分
	AppId              interface{} // 支付使用的APPID
	NotifyUrl          interface{} // 支付成功回调Url
	OpenApiId          interface{} // 使用的开放平台配置Id
	ChannelEdition     interface{} // 支付通道版本号
	TerminalIp         interface{} // 实时交易终端IP
	HidePaymentMethods interface{} // 隐藏支付方式，分号隔开;枚举： “INSTALMENT” “POSTPAID” “CARD” 在 GrabPay Checkout 流程中对用户隐藏指定的支付方式。如果未设置，GrabPay 会向用户显示所有符合条件的付款方式。但是请注意，您不能隐藏 GrabPay 钱包付款方式  注意：CARD 目前仅适用于泰国
	Verify             interface{} // codeVerify校验值
	Code               interface{} //
	Token              interface{} //
	AdditionalData     interface{} // 额外信息，JSON结构
	PaymentData        interface{} // 渠道支付接口返回核心参数，JSON结构
}
