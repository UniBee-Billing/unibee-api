// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentDao is the data access object for table payment.
type PaymentDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns PaymentColumns // columns contains all the column names of Table for convenient usage.
}

// PaymentColumns defines and stores column names for table payment.
type PaymentColumns struct {
	Id                 string // 主键id
	CompanyId          string // 公司Id
	MerchantId         string // 商户ID
	BizType            string // 业务类型。1-订单
	BizId              string // 业务id-即商户订单号
	CountryCode        string // 国家代码，指定发起交易的国家的两个字母的ISO 3166国家代码。目前支持SG、MY、PH、ID和TH
	Currency           string // 货币，“SGD” “MYR” “PHP” “IDR” “THB”\n与付款金额关联的货币。指定三个字母的ISO 4217货币代码
	MerchantOrderNo    string // 内部支付编号（系统生成唯一）
	PaymentFee         string // 支付金额
	RefundFee          string // 总共已退款金额
	BuyerPayFee        string // 买家实付金额
	ReceiptFee         string // 商户捕获金额（分）
	PayStatus          string // 支付状态。10-支付中，20-支付成功，30-支付取消
	AuthorizeStatus    string // 用户授权状态，0-未授权，1-已授权，2-已发起捕获
	ChannelId          string // 支付方式id,表oversea_pay_channel的id
	CaptureDelayHours  string // 延迟Capture时间
	ChannelPayId       string // 第三方支付平台支付预订单ID，支付接口返回
	ChannelTradeNo     string // 外部支付渠道订单号，支付成功回调返回
	CreateTime         string // 支付单创建时间
	CancelTime         string // 支付单取消时间
	PaidTime           string // 付款成功时间
	InvoiceTime        string // 入账成功时间
	GmtCreate          string // 创建时间
	GmtModify          string // 更新时间
	InvoiceStatus      string // 入账状态，未入账-0，入账中-1，完成入账-2，入账失败-3，已撤销入账-4
	InvoiceFee         string // 入账总金额。单位：分，invoice_total_fee + service_fee = payment_fee - refund_fee
	ServiceRate        string // 服务费比例，万分位，百分比[0，10000)，精度为0.01%，如3即为0.03%
	ServiceFee         string // 服务费。单位：分
	AppId              string // 支付使用的APPID
	NotifyUrl          string // 支付成功回调Url
	OpenApiId          string // 使用的开放平台配置Id
	ChannelEdition     string // 支付通道版本号
	TerminalIp         string // 实时交易终端IP
	HidePaymentMethods string // 隐藏支付方式，分号隔开;枚举： “INSTALMENT” “POSTPAID” “CARD”\n在 GrabPay Checkout 流程中对用户隐藏指定的支付方式。如果未设置，GrabPay 会向用户显示所有符合条件的付款方式。但是请注意，您不能隐藏 GrabPay 钱包付款方式\n\n注意：CARD 目前仅适用于泰国
	Verify             string // codeVerify校验值
	Code               string //
	Token              string //
	AdditionalData     string // 额外信息，JSON结构
	PaymentData        string // 渠道支付接口返回核心参数，JSON结构
}

// paymentColumns holds the columns for table payment.
var paymentColumns = PaymentColumns{
	Id:                 "id",
	CompanyId:          "company_id",
	MerchantId:         "merchant_id",
	BizType:            "biz_type",
	BizId:              "biz_id",
	CountryCode:        "country_code",
	Currency:           "currency",
	MerchantOrderNo:    "merchant_order_no",
	PaymentFee:         "payment_fee",
	RefundFee:          "refund_fee",
	BuyerPayFee:        "buyer_pay_fee",
	ReceiptFee:         "receipt_fee",
	PayStatus:          "pay_status",
	AuthorizeStatus:    "authorize_status",
	ChannelId:          "channel_id",
	CaptureDelayHours:  "capture_delay_hours",
	ChannelPayId:       "channel_pay_id",
	ChannelTradeNo:     "channel_trade_no",
	CreateTime:         "create_time",
	CancelTime:         "cancel_time",
	PaidTime:           "paid_time",
	InvoiceTime:        "invoice_time",
	GmtCreate:          "gmt_create",
	GmtModify:          "gmt_modify",
	InvoiceStatus:      "invoice_status",
	InvoiceFee:         "invoice_fee",
	ServiceRate:        "service_rate",
	ServiceFee:         "service_fee",
	AppId:              "app_id",
	NotifyUrl:          "notify_url",
	OpenApiId:          "open_api_id",
	ChannelEdition:     "channel_edition",
	TerminalIp:         "terminal_ip",
	HidePaymentMethods: "hide_payment_methods",
	Verify:             "verify",
	Code:               "code",
	Token:              "token",
	AdditionalData:     "additional_data",
	PaymentData:        "payment_data",
}

// NewPaymentDao creates and returns a new DAO object for table data access.
func NewPaymentDao() *PaymentDao {
	return &PaymentDao{
		group:   "oversea_pay",
		table:   "payment",
		columns: paymentColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PaymentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PaymentDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PaymentDao) Columns() PaymentColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PaymentDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PaymentDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PaymentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
