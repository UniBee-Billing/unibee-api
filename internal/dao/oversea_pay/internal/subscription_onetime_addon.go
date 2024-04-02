// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SubscriptionOnetimeAddonDao is the data access object for table subscription_onetime_addon.
type SubscriptionOnetimeAddonDao struct {
	table   string                          // table is the underlying table name of the DAO.
	group   string                          // group is the database configuration group name of current DAO.
	columns SubscriptionOnetimeAddonColumns // columns contains all the column names of Table for convenient usage.
}

// SubscriptionOnetimeAddonColumns defines and stores column names for table subscription_onetime_addon.
type SubscriptionOnetimeAddonColumns struct {
	Id             string // id
	GmtCreate      string // create_time
	GmtModify      string // modify_time
	SubscriptionId string // subscription_id
	AddonId        string // onetime addonId
	Quantity       string // quantity
	Status         string // status, 1-create, 2-paid, 3-cancel, 4-expired
	IsDeleted      string // 0-UnDeleted，1-Deleted
	CreateTime     string // create utc time
	PaymentId      string // paymentId
	MetaData       string // meta_data(json)
}

// subscriptionOnetimeAddonColumns holds the columns for table subscription_onetime_addon.
var subscriptionOnetimeAddonColumns = SubscriptionOnetimeAddonColumns{
	Id:             "id",
	GmtCreate:      "gmt_create",
	GmtModify:      "gmt_modify",
	SubscriptionId: "subscription_id",
	AddonId:        "addon_id",
	Quantity:       "quantity",
	Status:         "status",
	IsDeleted:      "is_deleted",
	CreateTime:     "create_time",
	PaymentId:      "payment_id",
	MetaData:       "meta_data",
}

// NewSubscriptionOnetimeAddonDao creates and returns a new DAO object for table data access.
func NewSubscriptionOnetimeAddonDao() *SubscriptionOnetimeAddonDao {
	return &SubscriptionOnetimeAddonDao{
		group:   "default",
		table:   "subscription_onetime_addon",
		columns: subscriptionOnetimeAddonColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SubscriptionOnetimeAddonDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SubscriptionOnetimeAddonDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SubscriptionOnetimeAddonDao) Columns() SubscriptionOnetimeAddonColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SubscriptionOnetimeAddonDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SubscriptionOnetimeAddonDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SubscriptionOnetimeAddonDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
