// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InvoiceDao is the data access object for table invoice.
type InvoiceDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns InvoiceColumns // columns contains all the column names of Table for convenient usage.
}

// InvoiceColumns defines and stores column names for table invoice.
type InvoiceColumns struct {
	Id                             string //
	MerchantId                     string // 商户Id
	UserId                         string // userId
	SubscriptionId                 string // 订阅id（内部编号）
	InvoiceId                      string // 发票ID（内部编号）
	GmtCreate                      string // 创建时间
	TotalAmount                    string // 金额,单位：分
	TaxAmount                      string // Tax金额,单位：分
	SubscriptionAmount             string // Sub金额,单位：分
	Currency                       string // 货币
	Lines                          string // lines json data
	ChannelId                      string // 支付渠道Id
	Status                         string // 订阅单状态，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
	SendStatus                     string // 邮件发送状态，0-No | 1- YES
	SendEmail                      string // email 发送地址，取自 UserAccount 表 email
	SendPdf                        string // pdf 文件地址
	Data                           string // 渠道额外参数，JSON格式
	GmtModify                      string // 修改时间
	IsDeleted                      string //
	Link                           string // invoice 链接（可用于支付）
	ChannelStatus                  string // 渠道最新状态，Stripe：https://stripe.com/docs/api/invoices/object
	ChannelPaymentId               string // 关联渠道 PaymentId
	ChannelUserId                  string // 渠道用户 Id
	ChannelInvoiceId               string // 关联渠道发票 Id
	ChannelInvoicePdf              string // 关联渠道发票 pdf
	TaxPercentage                  string // Tax税率，万分位，1000 表示 10%
	SendNote                       string // send_note
	SendTerms                      string // send_terms
	TotalAmountExcludingTax        string // 金额(不含税）,单位：分
	SubscriptionAmountExcludingTax string // Sub金额(不含税）,单位：分
	PeriodStart                    string // period_start
	PeriodEnd                      string // period_end
	PaymentId                      string // PaymentId
	RefundId                       string // refundId
	UniqueId                       string // 唯一键，stripe invoice 以同步为主，其他通道 invoice 实现方案不确定，使用自定义唯一键
}

// invoiceColumns holds the columns for table invoice.
var invoiceColumns = InvoiceColumns{
	Id:                             "id",
	MerchantId:                     "merchant_id",
	UserId:                         "user_id",
	SubscriptionId:                 "subscription_id",
	InvoiceId:                      "invoice_id",
	GmtCreate:                      "gmt_create",
	TotalAmount:                    "total_amount",
	TaxAmount:                      "tax_amount",
	SubscriptionAmount:             "subscription_amount",
	Currency:                       "currency",
	Lines:                          "lines",
	ChannelId:                      "channel_id",
	Status:                         "status",
	SendStatus:                     "send_status",
	SendEmail:                      "send_email",
	SendPdf:                        "send_pdf",
	Data:                           "data",
	GmtModify:                      "gmt_modify",
	IsDeleted:                      "is_deleted",
	Link:                           "link",
	ChannelStatus:                  "channel_status",
	ChannelPaymentId:               "channel_payment_id",
	ChannelUserId:                  "channel_user_id",
	ChannelInvoiceId:               "channel_invoice_id",
	ChannelInvoicePdf:              "channel_invoice_pdf",
	TaxPercentage:                  "tax_percentage",
	SendNote:                       "send_note",
	SendTerms:                      "send_terms",
	TotalAmountExcludingTax:        "total_amount_excluding_tax",
	SubscriptionAmountExcludingTax: "subscription_amount_excluding_tax",
	PeriodStart:                    "period_start",
	PeriodEnd:                      "period_end",
	PaymentId:                      "payment_id",
	RefundId:                       "refund_id",
	UniqueId:                       "unique_id",
}

// NewInvoiceDao creates and returns a new DAO object for table data access.
func NewInvoiceDao() *InvoiceDao {
	return &InvoiceDao{
		group:   "oversea_pay",
		table:   "invoice",
		columns: invoiceColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *InvoiceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *InvoiceDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *InvoiceDao) Columns() InvoiceColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *InvoiceDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *InvoiceDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *InvoiceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
