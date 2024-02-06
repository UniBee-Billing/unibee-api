// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SubscriptionTimelineDao is the data access object for table subscription_timeline.
type SubscriptionTimelineDao struct {
	table   string                      // table is the underlying table name of the DAO.
	group   string                      // group is the database configuration group name of current DAO.
	columns SubscriptionTimelineColumns // columns contains all the column names of Table for convenient usage.
}

// SubscriptionTimelineColumns defines and stores column names for table subscription_timeline.
type SubscriptionTimelineColumns struct {
	Id              string //
	MerchantId      string // merchant id
	UserId          string // userId
	SubscriptionId  string // subscription id
	PeriodStart     string // period_start
	PeriodEnd       string // period_end
	PeriodStartTime string // period start (datetime)
	PeriodEndTime   string // period end (datatime)
	GmtCreate       string // create time
	GmtModify       string // update time
	InvoiceId       string // invoice id
	UniqueId        string // unique id
	Currency        string // currency
	PlanId          string // PlanId
	Quantity        string // quantity
	AddonData       string // plan addon json data
	GatewayId       string // gateway_id
	IsDeleted       string // 0-UnDeleted，1-Deleted
	UniqueKey       string // unique key (deperated)
	CreateAt        string // create utc time
}

// subscriptionTimelineColumns holds the columns for table subscription_timeline.
var subscriptionTimelineColumns = SubscriptionTimelineColumns{
	Id:              "id",
	MerchantId:      "merchant_id",
	UserId:          "user_id",
	SubscriptionId:  "subscription_id",
	PeriodStart:     "period_start",
	PeriodEnd:       "period_end",
	PeriodStartTime: "period_start_time",
	PeriodEndTime:   "period_end_time",
	GmtCreate:       "gmt_create",
	GmtModify:       "gmt_modify",
	InvoiceId:       "invoice_id",
	UniqueId:        "unique_id",
	Currency:        "currency",
	PlanId:          "plan_id",
	Quantity:        "quantity",
	AddonData:       "addon_data",
	GatewayId:       "gateway_id",
	IsDeleted:       "is_deleted",
	UniqueKey:       "unique_key",
	CreateAt:        "create_at",
}

// NewSubscriptionTimelineDao creates and returns a new DAO object for table data access.
func NewSubscriptionTimelineDao() *SubscriptionTimelineDao {
	return &SubscriptionTimelineDao{
		group:   "oversea_pay",
		table:   "subscription_timeline",
		columns: subscriptionTimelineColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SubscriptionTimelineDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SubscriptionTimelineDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SubscriptionTimelineDao) Columns() SubscriptionTimelineColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SubscriptionTimelineDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SubscriptionTimelineDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SubscriptionTimelineDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
