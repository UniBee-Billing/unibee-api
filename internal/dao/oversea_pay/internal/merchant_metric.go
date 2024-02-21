// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantMetricDao is the data access object for table merchant_metric.
type MerchantMetricDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns MerchantMetricColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantMetricColumns defines and stores column names for table merchant_metric.
type MerchantMetricColumns struct {
	Id                  string // Id
	MerchantId          string // merchantId
	Code                string // code
	MetricName          string // metric name
	MetricDescription   string // metric description
	Type                string // 1-limit_metered，2-charge_metered(come later),3-charge_recurring(come later)
	AggregationType     string // 1-count，2-count unique, 3-latest, 4-max, 5-sum
	AggregationProperty string // aggregation property
	GmtCreate           string // create time
	GmtModify           string // update time
	IsDeleted           string // 0-UnDeleted，1-Deleted
	CreateTime          string // create utc time
}

// merchantMetricColumns holds the columns for table merchant_metric.
var merchantMetricColumns = MerchantMetricColumns{
	Id:                  "id",
	MerchantId:          "merchant_id",
	Code:                "code",
	MetricName:          "metric_name",
	MetricDescription:   "metric_description",
	Type:                "type",
	AggregationType:     "aggregation_type",
	AggregationProperty: "aggregation_property",
	GmtCreate:           "gmt_create",
	GmtModify:           "gmt_modify",
	IsDeleted:           "is_deleted",
	CreateTime:          "create_time",
}

// NewMerchantMetricDao creates and returns a new DAO object for table data access.
func NewMerchantMetricDao() *MerchantMetricDao {
	return &MerchantMetricDao{
		group:   "oversea_pay",
		table:   "merchant_metric",
		columns: merchantMetricColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantMetricDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantMetricDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantMetricDao) Columns() MerchantMetricColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantMetricDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantMetricDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantMetricDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
