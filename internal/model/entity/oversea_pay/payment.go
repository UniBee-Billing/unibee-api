// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Payment is the golang structure for table payment.
type Payment struct {
	Id                     int64       `json:"id"                     ` // 主键id
	CompanyId              int64       `json:"companyId"              ` // 公司Id
	MerchantId             int64       `json:"merchantId"             ` // 商户ID
	OpenApiId              int64       `json:"openApiId"              ` // 使用的开放平台配置Id
	UserId                 int64       `json:"userId"                 ` // user_id
	BizType                int         `json:"bizType"                ` // 业务类型。1-single payment, 3-invoice
	BizId                  string      `json:"bizId"                  ` // 业务id-即商户订单号
	TerminalIp             string      `json:"terminalIp"             ` // 实时交易终端IP
	CountryCode            string      `json:"countryCode"            ` // 国家代码，指定发起交易的国家的两个字母的ISO 3166国家代码。目前支持SG、MY、PH、ID和TH
	Currency               string      `json:"currency"               ` // 货币，“SGD” “MYR” “PHP” “IDR” “THB” 与付款金额关联的货币。指定三个字母的ISO 4217货币代码
	PaymentId              string      `json:"paymentId"              ` // 内部支付编号（系统生成唯一）
	PaymentFee             int64       `json:"paymentFee"             ` // 支付金额
	RefundFee              int64       `json:"refundFee"              ` // 总共已退款金额
	ReceiptFee             int64       `json:"receiptFee"             ` // 商户捕获金额（分）
	Status                 int         `json:"status"                 ` // 支付状态。10-支付中，20-支付成功，30-支付取消
	AuthorizeStatus        int         `json:"authorizeStatus"        ` // 用户授权状态，0-未授权，1-已授权，2-已发起捕获
	ChannelId              int64       `json:"channelId"              ` // 支付方式id,表oversea_pay_channel的id
	ChannelPaymentFee      int64       `json:"channelPaymentFee"      ` // 买家实付金额
	ChannelPaymentIntentId string      `json:"channelPaymentIntentId" ` // 第三方支付平台支付预订单ID，支付接口返回
	ChannelPaymentId       string      `json:"channelPaymentId"       ` // 外部支付渠道订单号，支付成功回调返回
	CaptureDelayHours      int         `json:"captureDelayHours"      ` // 延迟Capture时间
	CreateTime             *gtime.Time `json:"createTime"             ` // 支付单创建时间
	CancelTime             *gtime.Time `json:"cancelTime"             ` // 支付单取消时间
	PaidTime               *gtime.Time `json:"paidTime"               ` // 付款成功时间
	ChannelInvoiceId       string      `json:"channelInvoiceId"       ` // 渠道发票号
	InvoiceId              string      `json:"invoiceId"              ` // 发票号
	GmtCreate              *gtime.Time `json:"gmtCreate"              ` // 创建时间
	GmtModify              *gtime.Time `json:"gmtModify"              ` // 更新时间
	AppId                  string      `json:"appId"                  ` // 支付使用的APPID
	ReturnUrl              string      `json:"returnUrl"              ` // 支付成功回调Url
	ChannelEdition         string      `json:"channelEdition"         ` // 支付通道版本号
	HidePaymentMethods     string      `json:"hidePaymentMethods"     ` // 隐藏支付方式，分号隔开;枚举： “INSTALMENT” “POSTPAID” “CARD” 在 GrabPay Checkout 流程中对用户隐藏指定的支付方式。如果未设置，GrabPay 会向用户显示所有符合条件的付款方式。但是请注意，您不能隐藏 GrabPay 钱包付款方式  注意：CARD 目前仅适用于泰国
	Verify                 string      `json:"verify"                 ` // codeVerify校验值
	Code                   string      `json:"code"                   ` //
	Token                  string      `json:"token"                  ` //
	AdditionalData         string      `json:"additionalData"         ` // 额外信息，JSON结构
	PaymentData            string      `json:"paymentData"            ` // 渠道支付接口返回核心参数，JSON结构
}
