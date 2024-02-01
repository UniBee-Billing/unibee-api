// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EmailTemplateDao is the data access object for table email_template.
type EmailTemplateDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns EmailTemplateColumns // columns contains all the column names of Table for convenient usage.
}

// EmailTemplateColumns defines and stores column names for table email_template.
type EmailTemplateColumns struct {
	Id                 string //
	MerchantId         string //
	TemplateName       string //
	TemplateTitle      string //
	TemplateContent    string //
	TemplateAttachName string //
	GmtCreate          string // create time
	GmtModify          string // update time
	IsDeleted          string // 0-UnDeleted，1-Deleted
}

// emailTemplateColumns holds the columns for table email_template.
var emailTemplateColumns = EmailTemplateColumns{
	Id:                 "id",
	MerchantId:         "merchant_id",
	TemplateName:       "template_name",
	TemplateTitle:      "template_title",
	TemplateContent:    "template_content",
	TemplateAttachName: "template_attach_name",
	GmtCreate:          "gmt_create",
	GmtModify:          "gmt_modify",
	IsDeleted:          "is_deleted",
}

// NewEmailTemplateDao creates and returns a new DAO object for table data access.
func NewEmailTemplateDao() *EmailTemplateDao {
	return &EmailTemplateDao{
		group:   "oversea_pay",
		table:   "email_template",
		columns: emailTemplateColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *EmailTemplateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *EmailTemplateDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *EmailTemplateDao) Columns() EmailTemplateColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *EmailTemplateDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *EmailTemplateDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *EmailTemplateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
