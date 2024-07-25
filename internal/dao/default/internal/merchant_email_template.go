// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantEmailTemplateDao is the data access object for table merchant_email_template.
type MerchantEmailTemplateDao struct {
	table   string                       // table is the underlying table name of the DAO.
	group   string                       // group is the database configuration group name of current DAO.
	columns MerchantEmailTemplateColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantEmailTemplateColumns defines and stores column names for table merchant_email_template.
type MerchantEmailTemplateColumns struct {
	Id                 string //
	MerchantId         string //
	TemplateName       string //
	TemplateTitle      string //
	TemplateContent    string //
	TemplateAttachName string //
	GmtCreate          string // create time
	GmtModify          string // update time
	IsDeleted          string // 0-UnDeleted，1-Deleted
	CreateTime         string // create utc time
	Status             string // 0-Active,1-InActive
}

// merchantEmailTemplateColumns holds the columns for table merchant_email_template.
var merchantEmailTemplateColumns = MerchantEmailTemplateColumns{
	Id:                 "id",
	MerchantId:         "merchant_id",
	TemplateName:       "template_name",
	TemplateTitle:      "template_title",
	TemplateContent:    "template_content",
	TemplateAttachName: "template_attach_name",
	GmtCreate:          "gmt_create",
	GmtModify:          "gmt_modify",
	IsDeleted:          "is_deleted",
	CreateTime:         "create_time",
	Status:             "status",
}

// NewMerchantEmailTemplateDao creates and returns a new DAO object for table data access.
func NewMerchantEmailTemplateDao() *MerchantEmailTemplateDao {
	return &MerchantEmailTemplateDao{
		group:   "default",
		table:   "merchant_email_template",
		columns: merchantEmailTemplateColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantEmailTemplateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantEmailTemplateDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantEmailTemplateDao) Columns() MerchantEmailTemplateColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantEmailTemplateDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantEmailTemplateDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantEmailTemplateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
