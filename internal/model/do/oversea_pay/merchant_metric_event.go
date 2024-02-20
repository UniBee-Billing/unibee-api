// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantMetricEvent is the golang structure of table merchant_metric_event for DAO operations like Where/Data.
type MerchantMetricEvent struct {
	g.Meta                      `orm:"table:merchant_metric_event, do:true"`
	Id                          interface{} // Id
	MerchantId                  interface{} // merchantId
	MetricId                    interface{} // metric_id
	ExternalEventId             interface{} //
	UserId                      interface{} // metric_id
	AggregationPropertyInt      interface{} // aggregation property int, use for metric of max|sum type
	AggregationPropertyString   interface{} // aggregation property string, use for metric of count|count_unique type
	GmtCreate                   *gtime.Time // create time
	GmtModify                   *gtime.Time // update time
	IsDeleted                   interface{} // 0-UnDeleted，1-Deleted
	CreateTime                  interface{} // create utc time
	AggregationPropertyData     interface{} // aggregation property data (Json)
	AggregationPropertyUniqueId interface{} //
	SubscriptionIds             interface{} //
	SubscriptionPeriodStart     interface{} // matched subscription's current_period_start
	SubscriptionPeriodEnd       interface{} // matched subscription's current_period_end
}
