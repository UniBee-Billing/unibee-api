// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantMetric is the golang structure of table merchant_metric for DAO operations like Where/Data.
type MerchantMetric struct {
	g.Meta              `orm:"table:merchant_metric, do:true"`
	Id                  interface{} // Id
	MerchantId          interface{} // merchantId
	Code                interface{} // code
	MetricName          interface{} // metric name
	MetricDescription   interface{} // metric description
	Type                interface{} // 1-limit_metered，2-charge_metered(come later),3-charge_recurring(come later)
	AggregationType     interface{} // 1-count，2-count unique, 3-latest, 4-max, 5-sum
	AggregationProperty interface{} // aggregation property
	GmtCreate           *gtime.Time // create time
	GmtModify           *gtime.Time // update time
	IsDeleted           interface{} // 0-UnDeleted，1-Deleted
	CreateTime          interface{} // create utc time
}
