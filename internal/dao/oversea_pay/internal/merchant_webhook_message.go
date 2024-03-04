// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantWebhookMessageDao is the data access object for table merchant_webhook_message.
type MerchantWebhookMessageDao struct {
	table   string                        // table is the underlying table name of the DAO.
	group   string                        // group is the database configuration group name of current DAO.
	columns MerchantWebhookMessageColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantWebhookMessageColumns defines and stores column names for table merchant_webhook_message.
type MerchantWebhookMessageColumns struct {
	Id              string // id
	MerchantId      string // merchantId
	WebhookEvent    string // webhook_event
	Data            string // data(json)
	WebsocketStatus string // status  10-pending，20-success
	GmtCreate       string // create time
	GmtModify       string // update time
	CreateTime      string // create utc time
}

// merchantWebhookMessageColumns holds the columns for table merchant_webhook_message.
var merchantWebhookMessageColumns = MerchantWebhookMessageColumns{
	Id:              "id",
	MerchantId:      "merchant_id",
	WebhookEvent:    "webhook_event",
	Data:            "data",
	WebsocketStatus: "websocket_status",
	GmtCreate:       "gmt_create",
	GmtModify:       "gmt_modify",
	CreateTime:      "create_time",
}

// NewMerchantWebhookMessageDao creates and returns a new DAO object for table data access.
func NewMerchantWebhookMessageDao() *MerchantWebhookMessageDao {
	return &MerchantWebhookMessageDao{
		group:   "oversea_pay",
		table:   "merchant_webhook_message",
		columns: merchantWebhookMessageColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantWebhookMessageDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantWebhookMessageDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantWebhookMessageDao) Columns() MerchantWebhookMessageColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantWebhookMessageDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantWebhookMessageDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantWebhookMessageDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
