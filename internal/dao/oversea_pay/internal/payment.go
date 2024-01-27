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
	Id                     string // 主键id
	CompanyId              string // 公司Id
	MerchantId             string // 商户ID
	OpenApiId              string // 使用的开放平台配置Id
	UserId                 string // user_id
	SubscriptionId         string // 订阅id（内部编号）
	GmtCreate              string // 创建时间
	BizType                string // 业务类型。1-single payment, 3-subscription
	BizId                  string // 业务id-即商户订单号
	Currency               string // 货币，“SGD” “MYR” “PHP” “IDR” “THB”\n与付款金额关联的货币。指定三个字母的ISO 4217货币代码
	PaymentId              string // 内部支付编号（系统生成唯一）
	TotalAmount            string // 总计金额
	PaymentAmount          string // payment_amount
	BalanceAmount          string // balance_amount
	RefundAmount           string // 总共已退款金额
	Status                 string // 支付状态。10-支付中，20-支付成功，30-支付取消
	TerminalIp             string // 实时交易终端IP
	CountryCode            string // 国家代码，指定发起交易的国家的两个字母的ISO 3166国家代码。目前支持SG、MY、PH、ID和TH
	AuthorizeStatus        string // 用户授权状态，0-未授权，1-已授权，2-已发起捕获
	ChannelId              string // 支付方式id,表oversea_pay_channel的id
	ChannelPaymentIntentId string // 第三方支付平台支付预订单ID，支付接口返回
	ChannelPaymentId       string // 外部支付渠道订单号，支付成功回调返回
	CaptureDelayHours      string // 延迟Capture时间
	CreateTime             string // 支付单创建时间
	CancelTime             string // 支付单取消时间
	PaidTime               string // 付款成功时间
	InvoiceId              string // 发票号
	GmtModify              string // 更新时间
	AppId                  string // 支付使用的APPID
	ReturnUrl              string // 支付成功回调Url
	ChannelEdition         string // 支付通道版本号
	HidePaymentMethods     string // 隐藏支付方式，分号隔开;枚举： “INSTALMENT” “POSTPAID” “CARD”\n在 GrabPay Checkout 流程中对用户隐藏指定的支付方式。如果未设置，GrabPay 会向用户显示所有符合条件的付款方式。但是请注意，您不能隐藏 GrabPay 钱包付款方式\n\n注意：CARD 目前仅适用于泰国
	Verify                 string // codeVerify校验值
	Code                   string //
	Token                  string //
	AdditionalData         string // 额外信息，JSON结构
	PaymentData            string // 渠道支付接口返回核心参数，JSON结构
	UniqueId               string // 唯一键，以同步为逻辑加入使用自定义唯一键
	BalanceStart           string // balance_start
	BalanceEnd             string // balance_end
	InvoiceData            string //
}

// paymentColumns holds the columns for table payment.
var paymentColumns = PaymentColumns{
	Id:                     "id",
	CompanyId:              "company_id",
	MerchantId:             "merchant_id",
	OpenApiId:              "open_api_id",
	UserId:                 "user_id",
	SubscriptionId:         "subscription_id",
	GmtCreate:              "gmt_create",
	BizType:                "biz_type",
	BizId:                  "biz_id",
	Currency:               "currency",
	PaymentId:              "payment_id",
	TotalAmount:            "total_amount",
	PaymentAmount:          "payment_amount",
	BalanceAmount:          "balance_amount",
	RefundAmount:           "refund_amount",
	Status:                 "status",
	TerminalIp:             "terminal_ip",
	CountryCode:            "country_code",
	AuthorizeStatus:        "authorize_status",
	ChannelId:              "channel_id",
	ChannelPaymentIntentId: "channel_payment_intent_id",
	ChannelPaymentId:       "channel_payment_id",
	CaptureDelayHours:      "capture_delay_hours",
	CreateTime:             "create_time",
	CancelTime:             "cancel_time",
	PaidTime:               "paid_time",
	InvoiceId:              "invoice_id",
	GmtModify:              "gmt_modify",
	AppId:                  "app_id",
	ReturnUrl:              "return_url",
	ChannelEdition:         "channel_edition",
	HidePaymentMethods:     "hide_payment_methods",
	Verify:                 "verify",
	Code:                   "code",
	Token:                  "token",
	AdditionalData:         "additional_data",
	PaymentData:            "payment_data",
	UniqueId:               "unique_id",
	BalanceStart:           "balance_start",
	BalanceEnd:             "balance_end",
	InvoiceData:            "invoice_data",
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
