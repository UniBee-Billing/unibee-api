// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantBatchTaskExportChunksDao is the data access object for table merchant_batch_task_export_chunks.
type MerchantBatchTaskExportChunksDao struct {
	table   string                               // table is the underlying table name of the DAO.
	group   string                               // group is the database configuration group name of current DAO.
	columns MerchantBatchTaskExportChunksColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantBatchTaskExportChunksColumns defines and stores column names for table merchant_batch_task_export_chunks.
type MerchantBatchTaskExportChunksColumns struct {
	Id         string // id
	MerchantId string // merchant_id
	TaskId     string // task_id
	Page       string // page
	Content    string // content
	IsDeleted  string // 0-UnDeleted，1-Deleted
	GmtCreate  string // gmt_create
	GmtModify  string // update time
	CreateTime string // create utc time
}

// merchantBatchTaskExportChunksColumns holds the columns for table merchant_batch_task_export_chunks.
var merchantBatchTaskExportChunksColumns = MerchantBatchTaskExportChunksColumns{
	Id:         "id",
	MerchantId: "merchant_id",
	TaskId:     "task_id",
	Page:       "page",
	Content:    "content",
	IsDeleted:  "is_deleted",
	GmtCreate:  "gmt_create",
	GmtModify:  "gmt_modify",
	CreateTime: "create_time",
}

// NewMerchantBatchTaskExportChunksDao creates and returns a new DAO object for table data access.
func NewMerchantBatchTaskExportChunksDao() *MerchantBatchTaskExportChunksDao {
	return &MerchantBatchTaskExportChunksDao{
		group:   "default",
		table:   "merchant_batch_task_export_chunks",
		columns: merchantBatchTaskExportChunksColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantBatchTaskExportChunksDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantBatchTaskExportChunksDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantBatchTaskExportChunksDao) Columns() MerchantBatchTaskExportChunksColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantBatchTaskExportChunksDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantBatchTaskExportChunksDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantBatchTaskExportChunksDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
