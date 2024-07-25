// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantBatchExportTemplateDao is the data access object for table merchant_batch_export_template.
type MerchantBatchExportTemplateDao struct {
	table   string                             // table is the underlying table name of the DAO.
	group   string                             // group is the database configuration group name of current DAO.
	columns MerchantBatchExportTemplateColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantBatchExportTemplateColumns defines and stores column names for table merchant_batch_export_template.
type MerchantBatchExportTemplateColumns struct {
	Id            string // id
	MerchantId    string // merchant_id
	MemberId      string // member_id
	Name          string // name
	Task          string // task
	Format        string // format
	Payload       string // payload(json)
	ExportColumns string // export_columns(json)
	IsDeleted     string // 0-UnDeleted，1-Deleted
	GmtCreate     string // gmt_create
	GmtModify     string // update time
	CreateTime    string // create utc time
}

// merchantBatchExportTemplateColumns holds the columns for table merchant_batch_export_template.
var merchantBatchExportTemplateColumns = MerchantBatchExportTemplateColumns{
	Id:            "id",
	MerchantId:    "merchant_id",
	MemberId:      "member_id",
	Name:          "name",
	Task:          "task",
	Format:        "format",
	Payload:       "payload",
	ExportColumns: "export_columns",
	IsDeleted:     "is_deleted",
	GmtCreate:     "gmt_create",
	GmtModify:     "gmt_modify",
	CreateTime:    "create_time",
}

// NewMerchantBatchExportTemplateDao creates and returns a new DAO object for table data access.
func NewMerchantBatchExportTemplateDao() *MerchantBatchExportTemplateDao {
	return &MerchantBatchExportTemplateDao{
		group:   "default",
		table:   "merchant_batch_export_template",
		columns: merchantBatchExportTemplateColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantBatchExportTemplateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantBatchExportTemplateDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantBatchExportTemplateDao) Columns() MerchantBatchExportTemplateColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantBatchExportTemplateDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantBatchExportTemplateDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantBatchExportTemplateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
