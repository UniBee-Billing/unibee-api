// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantWebhookDao is the data access object for table merchant_webhook.
type MerchantWebhookDao struct {
	table   string                 // table is the underlying table name of the DAO.
	group   string                 // group is the database configuration group name of current DAO.
	columns MerchantWebhookColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantWebhookColumns defines and stores column names for table merchant_webhook.
type MerchantWebhookColumns struct {
	Id            string // id
	MerchantId    string // webhook url
	WebhookUrl    string // webhook url
	WebhookEvents string // webhook_events,split dot
	GmtCreate     string // create time
	GmtModify     string // update time
	CreateTime    string // create utc time
	IsDeleted     string // 0-UnDeleted，1-Deleted
}

// merchantWebhookColumns holds the columns for table merchant_webhook.
var merchantWebhookColumns = MerchantWebhookColumns{
	Id:            "id",
	MerchantId:    "merchant_id",
	WebhookUrl:    "webhook_url",
	WebhookEvents: "webhook_events",
	GmtCreate:     "gmt_create",
	GmtModify:     "gmt_modify",
	CreateTime:    "create_time",
	IsDeleted:     "is_deleted",
}

// NewMerchantWebhookDao creates and returns a new DAO object for table data access.
func NewMerchantWebhookDao() *MerchantWebhookDao {
	return &MerchantWebhookDao{
		group:   "default",
		table:   "merchant_webhook",
		columns: merchantWebhookColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantWebhookDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantWebhookDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantWebhookDao) Columns() MerchantWebhookColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantWebhookDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantWebhookDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantWebhookDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
