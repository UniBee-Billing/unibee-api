// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantMetricPlanLimitDao is the data access object for table merchant_metric_plan_limit.
type MerchantMetricPlanLimitDao struct {
	table   string                         // table is the underlying table name of the DAO.
	group   string                         // group is the database configuration group name of current DAO.
	columns MerchantMetricPlanLimitColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantMetricPlanLimitColumns defines and stores column names for table merchant_metric_plan_limit.
type MerchantMetricPlanLimitColumns struct {
	Id          string // Id
	MerchantId  string // merchantId
	MetricId    string // metricId
	PlanId      string // plan_id
	MetricLimit string // plan metric limit
	GmtCreate   string // create time
	GmtModify   string // update time
	IsDeleted   string // 0-UnDeleted，1-Deleted
	CreateTime  string // create utc time
}

// merchantMetricPlanLimitColumns holds the columns for table merchant_metric_plan_limit.
var merchantMetricPlanLimitColumns = MerchantMetricPlanLimitColumns{
	Id:          "id",
	MerchantId:  "merchant_id",
	MetricId:    "metric_id",
	PlanId:      "plan_id",
	MetricLimit: "metric_limit",
	GmtCreate:   "gmt_create",
	GmtModify:   "gmt_modify",
	IsDeleted:   "is_deleted",
	CreateTime:  "create_time",
}

// NewMerchantMetricPlanLimitDao creates and returns a new DAO object for table data access.
func NewMerchantMetricPlanLimitDao() *MerchantMetricPlanLimitDao {
	return &MerchantMetricPlanLimitDao{
		group:   "oversea_pay",
		table:   "merchant_metric_plan_limit",
		columns: merchantMetricPlanLimitColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantMetricPlanLimitDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantMetricPlanLimitDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantMetricPlanLimitDao) Columns() MerchantMetricPlanLimitColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantMetricPlanLimitDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantMetricPlanLimitDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantMetricPlanLimitDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}