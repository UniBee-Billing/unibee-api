// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CreditRefundDao is the data access object for table credit_refund.
type CreditRefundDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns CreditRefundColumns // columns contains all the column names of Table for convenient usage.
}

// CreditRefundColumns defines and stores column names for table credit_refund.
type CreditRefundColumns struct {
	Id                     string // Id
	UserId                 string // user_id
	CreditId               string // id of credit account
	Currency               string // currency
	InvoiceId              string // invoice_id
	CreditPaymentId        string // credit refund id
	CreditRefundId         string // credit refund id
	ExternalCreditRefundId string // external credit refund id
	RefundAmount           string // total refund amount,cent
	RefundTime             string // refund time
	Name                   string // recharge transaction title
	Description            string // recharge transaction description
	GmtCreate              string // create time
	GmtModify              string // update time
	CreateTime             string // create utc time
	MerchantId             string // merchant id
	AccountType            string // type of credit account, 1-main recharge account, 2-promo credit account
}

// creditRefundColumns holds the columns for table credit_refund.
var creditRefundColumns = CreditRefundColumns{
	Id:                     "id",
	UserId:                 "user_id",
	CreditId:               "credit_id",
	Currency:               "currency",
	InvoiceId:              "invoice_id",
	CreditPaymentId:        "credit_payment_id",
	CreditRefundId:         "credit_refund_id",
	ExternalCreditRefundId: "external_credit_refund_id",
	RefundAmount:           "refund_amount",
	RefundTime:             "refund_time",
	Name:                   "name",
	Description:            "description",
	GmtCreate:              "gmt_create",
	GmtModify:              "gmt_modify",
	CreateTime:             "create_time",
	MerchantId:             "merchant_id",
	AccountType:            "account_type",
}

// NewCreditRefundDao creates and returns a new DAO object for table data access.
func NewCreditRefundDao() *CreditRefundDao {
	return &CreditRefundDao{
		group:   "default",
		table:   "credit_refund",
		columns: creditRefundColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CreditRefundDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CreditRefundDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CreditRefundDao) Columns() CreditRefundColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CreditRefundDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CreditRefundDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CreditRefundDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
