// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantMetricEventDao is the data access object for table merchant_metric_event.
type MerchantMetricEventDao struct {
	table   string                     // table is the underlying table name of the DAO.
	group   string                     // group is the database configuration group name of current DAO.
	columns MerchantMetricEventColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantMetricEventColumns defines and stores column names for table merchant_metric_event.
type MerchantMetricEventColumns struct {
	Id                          string // Id
	MerchantId                  string // merchantId
	MetricId                    string // metric_id
	ExternalEventId             string // external_event_id, should be unique
	UserId                      string // user_id
	AggregationPropertyInt      string // aggregation property int, use for metric of max|sum type
	AggregationPropertyString   string // aggregation property string, use for metric of count|count_unique type
	GmtCreate                   string // create time
	GmtModify                   string // update time
	IsDeleted                   string // 0-UnDeleted，1-Deleted
	CreateTime                  string // create utc time
	AggregationPropertyData     string // aggregation property data (Json)
	AggregationPropertyUniqueId string //
	SubscriptionIds             string //
	SubscriptionPeriodStart     string // matched subscription's current_period_start
	SubscriptionPeriodEnd       string // matched subscription's current_period_end
	MetricLimit                 string //
	Used                        string //
	ChargeInvoiceId             string // charge invoice id
	ChargeData                  string // charge data
	ChargeStatus                string // 0-Uncharged，1-charged
}

// merchantMetricEventColumns holds the columns for table merchant_metric_event.
var merchantMetricEventColumns = MerchantMetricEventColumns{
	Id:                          "id",
	MerchantId:                  "merchant_id",
	MetricId:                    "metric_id",
	ExternalEventId:             "external_event_id",
	UserId:                      "user_id",
	AggregationPropertyInt:      "aggregation_property_int",
	AggregationPropertyString:   "aggregation_property_string",
	GmtCreate:                   "gmt_create",
	GmtModify:                   "gmt_modify",
	IsDeleted:                   "is_deleted",
	CreateTime:                  "create_time",
	AggregationPropertyData:     "aggregation_property_data",
	AggregationPropertyUniqueId: "aggregation_property_unique_id",
	SubscriptionIds:             "subscription_ids",
	SubscriptionPeriodStart:     "subscription_period_start",
	SubscriptionPeriodEnd:       "subscription_period_end",
	MetricLimit:                 "metric_limit",
	Used:                        "used",
	ChargeInvoiceId:             "charge_invoice_id",
	ChargeData:                  "charge_data",
	ChargeStatus:                "charge_status",
}

// NewMerchantMetricEventDao creates and returns a new DAO object for table data access.
func NewMerchantMetricEventDao() *MerchantMetricEventDao {
	return &MerchantMetricEventDao{
		group:   "default",
		table:   "merchant_metric_event",
		columns: merchantMetricEventColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantMetricEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantMetricEventDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantMetricEventDao) Columns() MerchantMetricEventColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantMetricEventDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantMetricEventDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantMetricEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
