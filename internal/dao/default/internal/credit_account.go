// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CreditAccountDao is the data access object for table credit_account.
type CreditAccountDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns CreditAccountColumns // columns contains all the column names of Table for convenient usage.
}

// CreditAccountColumns defines and stores column names for table credit_account.
type CreditAccountColumns struct {
	Id             string // Id
	UserId         string // user_id
	Type           string // type of credit account, 1-main recharge account, 2-promo credit account
	Currency       string // currency
	Amount         string // credit amount,cent
	GmtCreate      string // create time
	GmtModify      string // update time
	CreateTime     string // create utc time
	MerchantId     string // merchant id
	RechargeEnable string // 0-yes, 1-no
	PayoutEnable   string // 0-yes, 1-no
}

// creditAccountColumns holds the columns for table credit_account.
var creditAccountColumns = CreditAccountColumns{
	Id:             "id",
	UserId:         "user_id",
	Type:           "type",
	Currency:       "currency",
	Amount:         "amount",
	GmtCreate:      "gmt_create",
	GmtModify:      "gmt_modify",
	CreateTime:     "create_time",
	MerchantId:     "merchant_id",
	RechargeEnable: "recharge_enable",
	PayoutEnable:   "payout_enable",
}

// NewCreditAccountDao creates and returns a new DAO object for table data access.
func NewCreditAccountDao() *CreditAccountDao {
	return &CreditAccountDao{
		group:   "default",
		table:   "credit_account",
		columns: creditAccountColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CreditAccountDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CreditAccountDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CreditAccountDao) Columns() CreditAccountColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CreditAccountDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CreditAccountDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CreditAccountDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
