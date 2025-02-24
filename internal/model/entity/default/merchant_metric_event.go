// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantMetricEvent is the golang structure for table merchant_metric_event.
type MerchantMetricEvent struct {
	Id                          uint64      `json:"id"                          description:"Id"`                                                                     // Id
	MerchantId                  uint64      `json:"merchantId"                  description:"merchantId"`                                                             // merchantId
	MetricId                    uint64      `json:"metricId"                    description:"metric_id"`                                                              // metric_id
	ExternalEventId             string      `json:"externalEventId"             description:"external_event_id, should be unique"`                                    // external_event_id, should be unique
	UserId                      int64       `json:"userId"                      description:"user_id"`                                                                // user_id
	AggregationPropertyInt      uint64      `json:"aggregationPropertyInt"      description:"aggregation property int, use for metric of max|sum type"`               // aggregation property int, use for metric of max|sum type
	AggregationPropertyString   string      `json:"aggregationPropertyString"   description:"aggregation property string, use for metric of count|count_unique type"` // aggregation property string, use for metric of count|count_unique type
	GmtCreate                   *gtime.Time `json:"gmtCreate"                   description:"create time"`                                                            // create time
	GmtModify                   *gtime.Time `json:"gmtModify"                   description:"update time"`                                                            // update time
	IsDeleted                   int         `json:"isDeleted"                   description:"0-UnDeleted，1-Deleted"`                                                  // 0-UnDeleted，1-Deleted
	CreateTime                  int64       `json:"createTime"                  description:"create utc time"`                                                        // create utc time
	AggregationPropertyData     string      `json:"aggregationPropertyData"     description:"aggregation property data (Json)"`                                       // aggregation property data (Json)
	AggregationPropertyUniqueId string      `json:"aggregationPropertyUniqueId" description:""`                                                                       //
	SubscriptionIds             string      `json:"subscriptionIds"             description:""`                                                                       //
	SubscriptionPeriodStart     int64       `json:"subscriptionPeriodStart"     description:"matched subscription's current_period_start"`                            // matched subscription's current_period_start
	SubscriptionPeriodEnd       int64       `json:"subscriptionPeriodEnd"       description:"matched subscription's current_period_end"`                              // matched subscription's current_period_end
	MetricLimit                 uint64      `json:"metricLimit"                 description:""`                                                                       //
	Used                        uint64      `json:"used"                        description:""`                                                                       //
	ChargeInvoiceId             string      `json:"chargeInvoiceId"             description:"charge invoice id"`                                                      // charge invoice id
	ChargeData                  string      `json:"chargeData"                  description:"charge data"`                                                            // charge data
}
