// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EmailDefaultTemplateDao is the data access object for table email_default_template.
type EmailDefaultTemplateDao struct {
	table   string                      // table is the underlying table name of the DAO.
	group   string                      // group is the database configuration group name of current DAO.
	columns EmailDefaultTemplateColumns // columns contains all the column names of Table for convenient usage.
}

// EmailDefaultTemplateColumns defines and stores column names for table email_default_template.
type EmailDefaultTemplateColumns struct {
	Id                  string //
	TemplateName        string //
	TemplateDescription string //
	TemplateTitle       string //
	TemplateContent     string //
	TemplateAttachName  string //
	GmtCreate           string // create time
	GmtModify           string // update time
	IsDeleted           string // 0-UnDeleted，1-Deleted
	CreateTime          string // create utc time
}

// emailDefaultTemplateColumns holds the columns for table email_default_template.
var emailDefaultTemplateColumns = EmailDefaultTemplateColumns{
	Id:                  "id",
	TemplateName:        "template_name",
	TemplateDescription: "template_description",
	TemplateTitle:       "template_title",
	TemplateContent:     "template_content",
	TemplateAttachName:  "template_attach_name",
	GmtCreate:           "gmt_create",
	GmtModify:           "gmt_modify",
	IsDeleted:           "is_deleted",
	CreateTime:          "create_time",
}

// NewEmailDefaultTemplateDao creates and returns a new DAO object for table data access.
func NewEmailDefaultTemplateDao() *EmailDefaultTemplateDao {
	return &EmailDefaultTemplateDao{
		group:   "oversea_pay",
		table:   "email_default_template",
		columns: emailDefaultTemplateColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *EmailDefaultTemplateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *EmailDefaultTemplateDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *EmailDefaultTemplateDao) Columns() EmailDefaultTemplateColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *EmailDefaultTemplateDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *EmailDefaultTemplateDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *EmailDefaultTemplateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
