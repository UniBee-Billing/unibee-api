// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantMetric is the golang structure for table merchant_metric.
type MerchantMetric struct {
	Id                  int64       `json:"id"                  description:"Id"`                                                                          // Id
	MerchantId          int64       `json:"merchantId"          description:"merchantId"`                                                                  // merchantId
	Code                string      `json:"code"                description:"code"`                                                                        // code
	MetricName          string      `json:"metricName"          description:"metric name"`                                                                 // metric name
	MetricDescription   string      `json:"metricDescription"   description:"metric description"`                                                          // metric description
	Type                int         `json:"type"                description:"1-limit_metered，2-charge_metered(come later),3-charge_recurring(come later)"` // 1-limit_metered，2-charge_metered(come later),3-charge_recurring(come later)
	AggregationType     int         `json:"aggregationType"     description:"0-count，1-count unique, 2-latest, 3-max, 4-sum"`                              // 0-count，1-count unique, 2-latest, 3-max, 4-sum
	AggregationProperty string      `json:"aggregationProperty" description:"aggregation property"`                                                        // aggregation property
	GmtCreate           *gtime.Time `json:"gmtCreate"           description:"create time"`                                                                 // create time
	GmtModify           *gtime.Time `json:"gmtModify"           description:"update time"`                                                                 // update time
	IsDeleted           int         `json:"isDeleted"           description:"0-UnDeleted，1-Deleted"`                                                       // 0-UnDeleted，1-Deleted
	CreateTime          int64       `json:"createTime"          description:"create utc time"`                                                             // create utc time
}
