// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GatewayPlanDao is the data access object for table gateway_plan.
type GatewayPlanDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns GatewayPlanColumns // columns contains all the column names of Table for convenient usage.
}

// GatewayPlanColumns defines and stores column names for table gateway_plan.
type GatewayPlanColumns struct {
	Id                   string //
	GmtCreate            string // create time
	GmtModify            string // update time
	PlanId               string // PlanId
	GatewayId            string // gateway_id
	Status               string // 0-Init | 1-Create｜2-Active｜3-Inactive
	GatewayPlanId        string // gateway_plan_id
	GatewayProductId     string // gateway_product_id
	GatewayPlanStatus    string // gateway_plan_status
	GatewayProductStatus string // gateway_product_status
	IsDeleted            string // 0-UnDeleted，1-Deleted
	Data                 string // data(json)
	CreateTime           string // create utc time
}

// gatewayPlanColumns holds the columns for table gateway_plan.
var gatewayPlanColumns = GatewayPlanColumns{
	Id:                   "id",
	GmtCreate:            "gmt_create",
	GmtModify:            "gmt_modify",
	PlanId:               "plan_id",
	GatewayId:            "gateway_id",
	Status:               "status",
	GatewayPlanId:        "gateway_plan_id",
	GatewayProductId:     "gateway_product_id",
	GatewayPlanStatus:    "gateway_plan_status",
	GatewayProductStatus: "gateway_product_status",
	IsDeleted:            "is_deleted",
	Data:                 "data",
	CreateTime:           "create_time",
}

// NewGatewayPlanDao creates and returns a new DAO object for table data access.
func NewGatewayPlanDao() *GatewayPlanDao {
	return &GatewayPlanDao{
		group:   "oversea_pay",
		table:   "gateway_plan",
		columns: gatewayPlanColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *GatewayPlanDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *GatewayPlanDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *GatewayPlanDao) Columns() GatewayPlanColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *GatewayPlanDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *GatewayPlanDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *GatewayPlanDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
