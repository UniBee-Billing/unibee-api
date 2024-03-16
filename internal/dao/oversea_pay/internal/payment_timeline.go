// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentTimelineDao is the data access object for table payment_timeline.
type PaymentTimelineDao struct {
	table   string                 // table is the underlying table name of the DAO.
	group   string                 // group is the database configuration group name of current DAO.
	columns PaymentTimelineColumns // columns contains all the column names of Table for convenient usage.
}

// PaymentTimelineColumns defines and stores column names for table payment_timeline.
type PaymentTimelineColumns struct {
	Id             string //
	MerchantId     string // merchant id
	UserId         string // userId
	SubscriptionId string // subscription id
	InvoiceId      string // invoice id
	UniqueId       string // unique id
	Currency       string // currency
	TotalAmount    string // total amount
	GatewayId      string // gateway id
	GmtCreate      string // create time
	GmtModify      string // update time
	IsDeleted      string // 0-UnDeleted，1-Deleted
	PaymentId      string // PaymentId
	Status         string // 0-pending, 1-success, 2-failure
	TimelineType   string // 0-pay, 1-refund
	CreateTime     string // create utc time
}

// paymentTimelineColumns holds the columns for table payment_timeline.
var paymentTimelineColumns = PaymentTimelineColumns{
	Id:             "id",
	MerchantId:     "merchant_id",
	UserId:         "user_id",
	SubscriptionId: "subscription_id",
	InvoiceId:      "invoice_id",
	UniqueId:       "unique_id",
	Currency:       "currency",
	TotalAmount:    "total_amount",
	GatewayId:      "gateway_id",
	GmtCreate:      "gmt_create",
	GmtModify:      "gmt_modify",
	IsDeleted:      "is_deleted",
	PaymentId:      "payment_id",
	Status:         "status",
	TimelineType:   "timeline_type",
	CreateTime:     "create_time",
}

// NewPaymentTimelineDao creates and returns a new DAO object for table data access.
func NewPaymentTimelineDao() *PaymentTimelineDao {
	return &PaymentTimelineDao{
		group:   "default",
		table:   "payment_timeline",
		columns: paymentTimelineColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PaymentTimelineDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PaymentTimelineDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PaymentTimelineDao) Columns() PaymentTimelineColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PaymentTimelineDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PaymentTimelineDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PaymentTimelineDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
