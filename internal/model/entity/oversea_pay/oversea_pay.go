// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OverseaPay is the golang structure for table oversea_pay.
type OverseaPay struct {
	Id                 int64       `json:"id"                 ` // 主键id
	CompanyId          int64       `json:"companyId"          ` // 公司Id
	MerchantId         int64       `json:"merchantId"         ` // 商户ID
	BizType            int         `json:"bizType"            ` // 业务类型。1-订单
	BizId              string      `json:"bizId"              ` // 业务id-即商户订单号
	CountryCode        string      `json:"countryCode"        ` // 国家代码，指定发起交易的国家的两个字母的ISO 3166国家代码。目前支持SG、MY、PH、ID和TH
	Currency           string      `json:"currency"           ` // 货币，“SGD” “MYR” “PHP” “IDR” “THB” 与付款金额关联的货币。指定三个字母的ISO 4217货币代码
	MerchantOrderNo    string      `json:"merchantOrderNo"    ` // 内部支付编号（系统生成唯一）
	PaymentFee         int64       `json:"paymentFee"         ` // 支付金额
	RefundFee          int64       `json:"refundFee"          ` // 总共已退款金额
	BuyerPayFee        int64       `json:"buyerPayFee"        ` // 买家实付金额
	ReceiptFee         int64       `json:"receiptFee"         ` // 商户捕获金额（分）
	PayStatus          int         `json:"payStatus"          ` // 支付状态。10-支付中，20-支付成功，30-支付取消
	AuthorizeStatus    int         `json:"authorizeStatus"    ` // 用户授权状态，0-未授权，1-已授权，2-已发起捕获
	ChannelId          int64       `json:"channelId"          ` // 支付方式id,表oversea_pay_channel的id
	CaptureDelayHours  int         `json:"captureDelayHours"  ` // 延迟Capture时间
	ChannelPayId       string      `json:"channelPayId"       ` // 第三方支付平台支付预订单ID，支付接口返回
	ChannelTradeNo     string      `json:"channelTradeNo"     ` // 外部支付渠道订单号，支付成功回调返回
	CreateTime         *gtime.Time `json:"createTime"         ` // 支付单创建时间
	CancelTime         *gtime.Time `json:"cancelTime"         ` // 支付单取消时间
	PaidTime           *gtime.Time `json:"paidTime"           ` // 付款成功时间
	InvoiceTime        *gtime.Time `json:"invoiceTime"        ` // 入账成功时间
	GmtCreate          *gtime.Time `json:"gmtCreate"          ` // 创建时间
	GmtModify          *gtime.Time `json:"gmtModify"          ` // 更新时间
	InvoiceStatus      int         `json:"invoiceStatus"      ` // 入账状态，未入账-0，入账中-1，完成入账-2，入账失败-3，已撤销入账-4
	InvoiceFee         uint        `json:"invoiceFee"         ` // 入账总金额。单位：分，invoice_total_fee + service_fee = payment_fee - refund_fee
	ServiceRate        int64       `json:"serviceRate"        ` // 服务费比例，万分位，百分比[0，10000)，精度为0.01%，如3即为0.03%
	ServiceFee         int64       `json:"serviceFee"         ` // 服务费。单位：分
	AppId              string      `json:"appId"              ` // 支付使用的APPID
	NotifyUrl          string      `json:"notifyUrl"          ` // 支付成功回调Url
	OpenApiId          int64       `json:"openApiId"          ` // 使用的开放平台配置Id
	ChannelEdition     string      `json:"channelEdition"     ` // 支付通道版本号
	TerminalIp         string      `json:"terminalIp"         ` // 实时交易终端IP
	HidePaymentMethods string      `json:"hidePaymentMethods" ` // 隐藏支付方式，分号隔开;枚举： “INSTALMENT” “POSTPAID” “CARD” 在 GrabPay Checkout 流程中对用户隐藏指定的支付方式。如果未设置，GrabPay 会向用户显示所有符合条件的付款方式。但是请注意，您不能隐藏 GrabPay 钱包付款方式  注意：CARD 目前仅适用于泰国
	Verify             string      `json:"verify"             ` // codeVerify校验值
	Code               string      `json:"code"               ` //
	Token              string      `json:"token"              ` //
	AdditionalData     string      `json:"additionalData"     ` // 额外信息，JSON结构
	PaymentData        string      `json:"paymentData"        ` // 渠道支付接口返回核心参数，JSON结构
}
