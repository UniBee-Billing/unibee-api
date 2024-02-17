// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantWebhookLogDao is the data access object for table merchant_webhook_log.
type MerchantWebhookLogDao struct {
	table   string                    // table is the underlying table name of the DAO.
	group   string                    // group is the database configuration group name of current DAO.
	columns MerchantWebhookLogColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantWebhookLogColumns defines and stores column names for table merchant_webhook_log.
type MerchantWebhookLogColumns struct {
	Id           string // id
	MerchantId   string // webhook url
	WebhookUrl   string // webhook url
	WebhookEvent string // webhook_event
	RequestId    string // request_id
	Body         string // body(json)
	Response     string // response
	Mamo         string // mamo
	GmtCreate    string // create time
	GmtModify    string // update time
	CreateTime   string // create utc time
}

// merchantWebhookLogColumns holds the columns for table merchant_webhook_log.
var merchantWebhookLogColumns = MerchantWebhookLogColumns{
	Id:           "id",
	MerchantId:   "merchant_id",
	WebhookUrl:   "webhook_url",
	WebhookEvent: "webhook_event",
	RequestId:    "request_id",
	Body:         "body",
	Response:     "response",
	Mamo:         "mamo",
	GmtCreate:    "gmt_create",
	GmtModify:    "gmt_modify",
	CreateTime:   "create_time",
}

// NewMerchantWebhookLogDao creates and returns a new DAO object for table data access.
func NewMerchantWebhookLogDao() *MerchantWebhookLogDao {
	return &MerchantWebhookLogDao{
		group:   "oversea_pay",
		table:   "merchant_webhook_log",
		columns: merchantWebhookLogColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantWebhookLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantWebhookLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantWebhookLogDao) Columns() MerchantWebhookLogColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantWebhookLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantWebhookLogDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantWebhookLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
