// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantBatchTaskDao is the data access object for table merchant_batch_task.
type MerchantBatchTaskDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns MerchantBatchTaskColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantBatchTaskColumns defines and stores column names for table merchant_batch_task.
type MerchantBatchTaskColumns struct {
	Id            string // id
	MerchantId    string // merchant_id
	MemberId      string // member_id
	ModuleName    string // module_name
	TaskName      string // task_name
	SourceFrom    string // source_from
	Payload       string // payload(json)
	DownloadUrl   string // download_file_url
	Status        string // Status。0-Pending，1-Processing，2-Success，3-Failure
	StartTime     string // task_start_time
	FinishTime    string // task_finish_time
	TaskCost      string // task cost time(second)
	FailReason    string // reason of failure
	GmtCreate     string // gmt_create
	TaskType      string // type，0-download，1-upload
	SuccessCount  string // success_count
	UploadFileUrl string // the file url of upload type task
	GmtModify     string // update time
}

// merchantBatchTaskColumns holds the columns for table merchant_batch_task.
var merchantBatchTaskColumns = MerchantBatchTaskColumns{
	Id:            "id",
	MerchantId:    "merchant_id",
	MemberId:      "member_id",
	ModuleName:    "module_name",
	TaskName:      "task_name",
	SourceFrom:    "source_from",
	Payload:       "payload",
	DownloadUrl:   "download_url",
	Status:        "status",
	StartTime:     "start_time",
	FinishTime:    "finish_time",
	TaskCost:      "task_cost",
	FailReason:    "fail_reason",
	GmtCreate:     "gmt_create",
	TaskType:      "task_type",
	SuccessCount:  "success_count",
	UploadFileUrl: "upload_file_url",
	GmtModify:     "gmt_modify",
}

// NewMerchantBatchTaskDao creates and returns a new DAO object for table data access.
func NewMerchantBatchTaskDao() *MerchantBatchTaskDao {
	return &MerchantBatchTaskDao{
		group:   "default",
		table:   "merchant_batch_task",
		columns: merchantBatchTaskColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantBatchTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantBatchTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantBatchTaskDao) Columns() MerchantBatchTaskColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantBatchTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantBatchTaskDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantBatchTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
