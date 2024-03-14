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
	Id                     string // id
	CompanyId              string // company id
	MerchantId             string // merchant id
	OpenApiId              string // open api id
	UserId                 string // user_id
	SubscriptionId         string // subscription id
	GmtCreate              string // create time
	GmtModify              string // update time
	BizType                string // biz_type 1-onetime payment, 3-subscription
	ExternalPaymentId      string // external_payment_id
	Currency               string // currency，“SGD” “MYR” “PHP” “IDR” “THB”
	PaymentId              string // payment id
	TotalAmount            string // total amount
	PaymentAmount          string // payment_amount
	BalanceAmount          string // balance_amount
	RefundAmount           string // total refund amount
	Status                 string // status  10-pending，20-success，30-failure, 40-cancel
	TerminalIp             string // client ip
	CountryCode            string // country code
	AuthorizeStatus        string // authorize status，0-waiting authorize，1-authorized，2-authorized_request
	AuthorizeReason        string //
	GatewayId              string // gateway_id
	GatewayPaymentIntentId string // gateway_payment_intent_id
	GatewayPaymentId       string // gateway_payment_id
	CaptureDelayHours      string // capture_delay_hours
	CreateTime             string // create time, utc time
	CancelTime             string // cancel time, utc time
	PaidTime               string // paid time, utc time
	InvoiceId              string // invoice id
	AppId                  string // app id
	ReturnUrl              string // return url
	GatewayEdition         string // gateway edition
	HidePaymentMethods     string // hide_payment_methods
	Verify                 string // codeVerify
	Code                   string //
	Token                  string //
	MetaData               string // meta_data (json)
	Automatic              string // 0-no,1-yes
	FailureReason          string //
	BillingReason          string //
	Link                   string //
	PaymentData            string // payment create data (json)
	UniqueId               string // unique id
	BalanceStart           string // balance_start, utc time
	BalanceEnd             string // balance_end, utc time
	InvoiceData            string //
	GatewayPaymentMethod   string //
	GasPayer               string // who pay the gas, merchant|user
	ExpireTime             string // expire time, utc time
	GatewayLink            string //
	CryptoAmount           string // crypto_amount, cent
	CryptoCurrency         string // crypto_currency
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
	GmtModify:              "gmt_modify",
	BizType:                "biz_type",
	ExternalPaymentId:      "external_payment_id",
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
	AuthorizeReason:        "authorize_reason",
	GatewayId:              "gateway_id",
	GatewayPaymentIntentId: "gateway_payment_intent_id",
	GatewayPaymentId:       "gateway_payment_id",
	CaptureDelayHours:      "capture_delay_hours",
	CreateTime:             "create_time",
	CancelTime:             "cancel_time",
	PaidTime:               "paid_time",
	InvoiceId:              "invoice_id",
	AppId:                  "app_id",
	ReturnUrl:              "return_url",
	GatewayEdition:         "gateway_edition",
	HidePaymentMethods:     "hide_payment_methods",
	Verify:                 "verify",
	Code:                   "code",
	Token:                  "token",
	MetaData:               "meta_data",
	Automatic:              "automatic",
	FailureReason:          "failure_reason",
	BillingReason:          "billing_reason",
	Link:                   "link",
	PaymentData:            "payment_data",
	UniqueId:               "unique_id",
	BalanceStart:           "balance_start",
	BalanceEnd:             "balance_end",
	InvoiceData:            "invoice_data",
	GatewayPaymentMethod:   "gateway_payment_method",
	GasPayer:               "gas_payer",
	ExpireTime:             "expire_time",
	GatewayLink:            "gateway_link",
	CryptoAmount:           "crypto_amount",
	CryptoCurrency:         "crypto_currency",
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
