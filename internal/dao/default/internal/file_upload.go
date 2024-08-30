// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FileUploadDao is the data access object for table file_upload.
type FileUploadDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns FileUploadColumns // columns contains all the column names of Table for convenient usage.
}

// FileUploadColumns defines and stores column names for table file_upload.
type FileUploadColumns struct {
	Id         string //
	UserId     string //
	Url        string //
	FileName   string //
	Tag        string //
	GmtCreate  string // create time
	GmtModify  string //
	IsDeleted  string // 0-UnDeleted，1-Deleted
	CreateTime string // create utc time
	Data       string //
}

// fileUploadColumns holds the columns for table file_upload.
var fileUploadColumns = FileUploadColumns{
	Id:         "id",
	UserId:     "user_id",
	Url:        "url",
	FileName:   "file_name",
	Tag:        "tag",
	GmtCreate:  "gmt_create",
	GmtModify:  "gmt_modify",
	IsDeleted:  "is_deleted",
	CreateTime: "create_time",
	Data:       "data",
}

// NewFileUploadDao creates and returns a new DAO object for table data access.
func NewFileUploadDao() *FileUploadDao {
	return &FileUploadDao{
		group:   "default",
		table:   "file_upload",
		columns: fileUploadColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *FileUploadDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *FileUploadDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *FileUploadDao) Columns() FileUploadColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *FileUploadDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *FileUploadDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *FileUploadDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
