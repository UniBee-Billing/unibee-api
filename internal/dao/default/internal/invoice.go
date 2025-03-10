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
	MetaData                       string // meta_data(json)
	BizType                        string // biz type from payment 1-single payment, 3-subscription
	MerchantId                     string // merchant_id
	UserId                         string // userId
	SubscriptionId                 string // subscription_id
	InvoiceId                      string // invoice_id
	InvoiceName                    string // invoice name
	UniqueId                       string // unique_id
	GmtCreate                      string // create time
	GmtModify                      string // update time
	TotalAmount                    string // total amount, cent
	TaxAmount                      string // tax amount,cent
	SubscriptionAmount             string // sub amount,cent
	Currency                       string // currency
	Lines                          string // lines( json)
	PaymentId                      string // paymentId
	GatewayId                      string // gateway_id
	Status                         string // status，0-Init | 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
	SendStatus                     string // email send status，0-No | 1- YES| 2-Unnecessary
	SendEmail                      string // email
	SendPdf                        string // pdf link
	IsDeleted                      string // 0-UnDeleted，1-Deleted
	Link                           string // invoice link
	PaymentLink                    string // invoice payment link
	GatewayStatus                  string //
	GatewayInvoiceId               string //
	GatewayPaymentId               string //
	GatewayInvoicePdf              string //
	TaxPercentage                  string // TaxPercentage，1000 = 10%
	SendNote                       string // send_note
	SendTerms                      string // send_terms
	TotalAmountExcludingTax        string //
	SubscriptionAmountExcludingTax string //
	PeriodStart                    string // period_start, utc time
	PeriodEnd                      string // period_end utc time
	TrialEnd                       string // trial_end, utc time
	PeriodStartTime                string //
	PeriodEndTime                  string //
	RefundId                       string // refundId
	Data                           string // data (json)
	CreateTime                     string // create utc time
	CryptoAmount                   string // crypto_amount, cent
	CryptoCurrency                 string // crypto_currency
	FinishTime                     string // utc time of enter process
	DayUtilDue                     string // day util due after process
	LastTrackTime                  string // last process invoice track time
	DiscountCode                   string // discount_code
	DiscountAmount                 string // discount amount, cent
	CountryCode                    string //
	ProductName                    string // product name
	GasPayer                       string // gas_payer
	GatewayPaymentMethod           string // gateway_payment_method
	BillingCycleAnchor             string // billing_cycle_anchor
	CreateFrom                     string // create from
	VatNumber                      string //
	PromoCreditDiscountAmount      string // promo credit discount amount
	PartialCreditPaidAmount        string // partial credit paid amount
	MetricCharge                   string // invoice metric charge data
}

// invoiceColumns holds the columns for table invoice.
var invoiceColumns = InvoiceColumns{
	Id:                             "id",
	MetaData:                       "meta_data",
	BizType:                        "biz_type",
	MerchantId:                     "merchant_id",
	UserId:                         "user_id",
	SubscriptionId:                 "subscription_id",
	InvoiceId:                      "invoice_id",
	InvoiceName:                    "invoice_name",
	UniqueId:                       "unique_id",
	GmtCreate:                      "gmt_create",
	GmtModify:                      "gmt_modify",
	TotalAmount:                    "total_amount",
	TaxAmount:                      "tax_amount",
	SubscriptionAmount:             "subscription_amount",
	Currency:                       "currency",
	Lines:                          "lines",
	PaymentId:                      "payment_id",
	GatewayId:                      "gateway_id",
	Status:                         "status",
	SendStatus:                     "send_status",
	SendEmail:                      "send_email",
	SendPdf:                        "send_pdf",
	IsDeleted:                      "is_deleted",
	Link:                           "link",
	PaymentLink:                    "payment_link",
	GatewayStatus:                  "gateway_status",
	GatewayInvoiceId:               "gateway_invoice_id",
	GatewayPaymentId:               "gateway_payment_id",
	GatewayInvoicePdf:              "gateway_invoice_pdf",
	TaxPercentage:                  "tax_percentage",
	SendNote:                       "send_note",
	SendTerms:                      "send_terms",
	TotalAmountExcludingTax:        "total_amount_excluding_tax",
	SubscriptionAmountExcludingTax: "subscription_amount_excluding_tax",
	PeriodStart:                    "period_start",
	PeriodEnd:                      "period_end",
	TrialEnd:                       "trial_end",
	PeriodStartTime:                "period_start_time",
	PeriodEndTime:                  "period_end_time",
	RefundId:                       "refund_id",
	Data:                           "data",
	CreateTime:                     "create_time",
	CryptoAmount:                   "crypto_amount",
	CryptoCurrency:                 "crypto_currency",
	FinishTime:                     "finish_time",
	DayUtilDue:                     "day_util_due",
	LastTrackTime:                  "last_track_time",
	DiscountCode:                   "discount_code",
	DiscountAmount:                 "discount_amount",
	CountryCode:                    "country_code",
	ProductName:                    "product_name",
	GasPayer:                       "gas_payer",
	GatewayPaymentMethod:           "gateway_payment_method",
	BillingCycleAnchor:             "billing_cycle_anchor",
	CreateFrom:                     "create_from",
	VatNumber:                      "vat_number",
	PromoCreditDiscountAmount:      "promo_credit_discount_amount",
	PartialCreditPaidAmount:        "partial_credit_paid_amount",
	MetricCharge:                   "metric_charge",
}

// NewInvoiceDao creates and returns a new DAO object for table data access.
func NewInvoiceDao() *InvoiceDao {
	return &InvoiceDao{
		group:   "default",
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
