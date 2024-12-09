// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CreditTransactionDao is the data access object for table credit_transaction.
type CreditTransactionDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns CreditTransactionColumns // columns contains all the column names of Table for convenient usage.
}

// CreditTransactionColumns defines and stores column names for table credit_transaction.
type CreditTransactionColumns struct {
	Id                 string // Id
	UserId             string // user_id
	CreditId           string // id of credit account
	Currency           string // currency
	TransactionId      string // unique id for timeline
	TransactionType    string // transaction type。1-recharge income，2-payment out，3-refund income，4-withdraw out，5-withdraw failed income, 6-admin change，7-recharge refund out
	CreditAmountAfter  string // the credit amount after transaction,cent
	CreditAmountBefore string // the credit amount before transaction,cent
	DeltaAmount        string // delta amount,cent
	BizId              string // bisness id
	Name               string // recharge transaction title
	Description        string // recharge transaction description
	GmtCreate          string // create time
	GmtModify          string // update time
	CreateTime         string // create utc time
	MerchantId         string // merchant id
	InvoiceId          string // invoice_id
	AccountType        string // type of credit account, 1-main recharge account, 2-promo credit account
	AdminMemberId      string // admin_member_id
}

// creditTransactionColumns holds the columns for table credit_transaction.
var creditTransactionColumns = CreditTransactionColumns{
	Id:                 "id",
	UserId:             "user_id",
	CreditId:           "credit_id",
	Currency:           "currency",
	TransactionId:      "transaction_id",
	TransactionType:    "transaction_type",
	CreditAmountAfter:  "credit_amount_after",
	CreditAmountBefore: "credit_amount_before",
	DeltaAmount:        "delta_amount",
	BizId:              "biz_id",
	Name:               "name",
	Description:        "description",
	GmtCreate:          "gmt_create",
	GmtModify:          "gmt_modify",
	CreateTime:         "create_time",
	MerchantId:         "merchant_id",
	InvoiceId:          "invoice_id",
	AccountType:        "account_type",
	AdminMemberId:      "admin_member_id",
}

// NewCreditTransactionDao creates and returns a new DAO object for table data access.
func NewCreditTransactionDao() *CreditTransactionDao {
	return &CreditTransactionDao{
		group:   "default",
		table:   "credit_transaction",
		columns: creditTransactionColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CreditTransactionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CreditTransactionDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CreditTransactionDao) Columns() CreditTransactionColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CreditTransactionDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CreditTransactionDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CreditTransactionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
