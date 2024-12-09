// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CreditPaymentDao is the data access object for table credit_payment.
type CreditPaymentDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns CreditPaymentColumns // columns contains all the column names of Table for convenient usage.
}

// CreditPaymentColumns defines and stores column names for table credit_payment.
type CreditPaymentColumns struct {
	Id                      string // Id
	UserId                  string // user_id
	CreditId                string // id of credit account
	Currency                string // currency
	CreditPaymentId         string // credit payment id
	ExternalCreditPaymentId string // external credit payment id
	TotalAmount             string // total amount,cent
	PaidTime                string // paid time
	Name                    string // recharge transaction title
	Description             string // recharge transaction description
	GmtCreate               string // create time
	GmtModify               string // update time
	CreateTime              string // create utc time
	MerchantId              string // merchant id
	InvoiceId               string // invoice_id
	TotalRefundAmount       string // total amount,cent
	ExchangeRate            string //
	PaidCurrencyAmount      string //
	AccountType             string // type of credit account, 1-main recharge account, 2-promo credit account
}

// creditPaymentColumns holds the columns for table credit_payment.
var creditPaymentColumns = CreditPaymentColumns{
	Id:                      "id",
	UserId:                  "user_id",
	CreditId:                "credit_id",
	Currency:                "currency",
	CreditPaymentId:         "credit_payment_id",
	ExternalCreditPaymentId: "external_credit_payment_id",
	TotalAmount:             "total_amount",
	PaidTime:                "paid_time",
	Name:                    "name",
	Description:             "description",
	GmtCreate:               "gmt_create",
	GmtModify:               "gmt_modify",
	CreateTime:              "create_time",
	MerchantId:              "merchant_id",
	InvoiceId:               "invoice_id",
	TotalRefundAmount:       "total_refund_amount",
	ExchangeRate:            "exchange_rate",
	PaidCurrencyAmount:      "paid_currency_amount",
	AccountType:             "account_type",
}

// NewCreditPaymentDao creates and returns a new DAO object for table data access.
func NewCreditPaymentDao() *CreditPaymentDao {
	return &CreditPaymentDao{
		group:   "default",
		table:   "credit_payment",
		columns: creditPaymentColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CreditPaymentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CreditPaymentDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CreditPaymentDao) Columns() CreditPaymentColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CreditPaymentDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CreditPaymentDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CreditPaymentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
