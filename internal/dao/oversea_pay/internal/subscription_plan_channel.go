// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SubscriptionPlanChannelDao is the data access object for table subscription_plan_channel.
type SubscriptionPlanChannelDao struct {
	table   string                         // table is the underlying table name of the DAO.
	group   string                         // group is the database configuration group name of current DAO.
	columns SubscriptionPlanChannelColumns // columns contains all the column names of Table for convenient usage.
}

// SubscriptionPlanChannelColumns defines and stores column names for table subscription_plan_channel.
type SubscriptionPlanChannelColumns struct {
	Id               string //
	GmtCreate        string // 创建时间
	GmtModify        string // 修改时间
	PlanId           string // 计划ID
	ChannelId        string // 支付渠道Id
	ChannelPlanId    string // 支付渠道plan_Id
	ChannelProductId string // 支付渠道product_Id
	Data             string // 渠道额外参数，JSON格式
	IsDeleted        string //
}

// subscriptionPlanChannelColumns holds the columns for table subscription_plan_channel.
var subscriptionPlanChannelColumns = SubscriptionPlanChannelColumns{
	Id:               "id",
	GmtCreate:        "gmt_create",
	GmtModify:        "gmt_modify",
	PlanId:           "plan_id",
	ChannelId:        "channel_id",
	ChannelPlanId:    "channel_plan_id",
	ChannelProductId: "channel_product_id",
	Data:             "data",
	IsDeleted:        "is_deleted",
}

// NewSubscriptionPlanChannelDao creates and returns a new DAO object for table data access.
func NewSubscriptionPlanChannelDao() *SubscriptionPlanChannelDao {
	return &SubscriptionPlanChannelDao{
		group:   "oversea_pay",
		table:   "subscription_plan_channel",
		columns: subscriptionPlanChannelColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SubscriptionPlanChannelDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SubscriptionPlanChannelDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SubscriptionPlanChannelDao) Columns() SubscriptionPlanChannelColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SubscriptionPlanChannelDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SubscriptionPlanChannelDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SubscriptionPlanChannelDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
