// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantEmailHistoryDao is the data access object for table merchant_email_history.
type MerchantEmailHistoryDao struct {
	table   string                      // table is the underlying table name of the DAO.
	group   string                      // group is the database configuration group name of current DAO.
	columns MerchantEmailHistoryColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantEmailHistoryColumns defines and stores column names for table merchant_email_history.
type MerchantEmailHistoryColumns struct {
	Id         string //
	MerchantId string //
	Email      string //
	Title      string //
	Content    string //
	AttachFile string //
	GmtCreate  string // create time
	GmtModify  string // update time
	Response   string //
	CreateTime string // create utc time
}

// merchantEmailHistoryColumns holds the columns for table merchant_email_history.
var merchantEmailHistoryColumns = MerchantEmailHistoryColumns{
	Id:         "id",
	MerchantId: "merchant_id",
	Email:      "email",
	Title:      "title",
	Content:    "content",
	AttachFile: "attach_file",
	GmtCreate:  "gmt_create",
	GmtModify:  "gmt_modify",
	Response:   "response",
	CreateTime: "create_time",
}

// NewMerchantEmailHistoryDao creates and returns a new DAO object for table data access.
func NewMerchantEmailHistoryDao() *MerchantEmailHistoryDao {
	return &MerchantEmailHistoryDao{
		group:   "oversea_pay",
		table:   "merchant_email_history",
		columns: merchantEmailHistoryColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantEmailHistoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantEmailHistoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantEmailHistoryDao) Columns() MerchantEmailHistoryColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantEmailHistoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantEmailHistoryDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantEmailHistoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
