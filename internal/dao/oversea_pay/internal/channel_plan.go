// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ChannelPlanDao is the data access object for table channel_plan.
type ChannelPlanDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns ChannelPlanColumns // columns contains all the column names of Table for convenient usage.
}

// ChannelPlanColumns defines and stores column names for table channel_plan.
type ChannelPlanColumns struct {
	Id                   string //
	GmtCreate            string // create time
	GmtModify            string // update time
	PlanId               string // PlanId
	ChannelId            string // channel_id
	Status               string // 0-Init | 1-Create｜2-Active｜3-Inactive
	ChannelPlanId        string // channel_plan_Id
	ChannelProductId     string // channel_product_Id
	ChannelPlanStatus    string // channel_plan_status
	ChannelProductStatus string // channel_product_status
	IsDeleted            string // 0-UnDeleted，1-Deleted
	Data                 string // data(json)
}

// channelPlanColumns holds the columns for table channel_plan.
var channelPlanColumns = ChannelPlanColumns{
	Id:                   "id",
	GmtCreate:            "gmt_create",
	GmtModify:            "gmt_modify",
	PlanId:               "plan_id",
	ChannelId:            "channel_id",
	Status:               "status",
	ChannelPlanId:        "channel_plan_id",
	ChannelProductId:     "channel_product_id",
	ChannelPlanStatus:    "channel_plan_status",
	ChannelProductStatus: "channel_product_status",
	IsDeleted:            "is_deleted",
	Data:                 "data",
}

// NewChannelPlanDao creates and returns a new DAO object for table data access.
func NewChannelPlanDao() *ChannelPlanDao {
	return &ChannelPlanDao{
		group:   "oversea_pay",
		table:   "channel_plan",
		columns: channelPlanColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ChannelPlanDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ChannelPlanDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ChannelPlanDao) Columns() ChannelPlanColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ChannelPlanDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ChannelPlanDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ChannelPlanDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
