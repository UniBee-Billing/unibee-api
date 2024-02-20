// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantMetricPlanLimit is the golang structure for table merchant_metric_plan_limit.
type MerchantMetricPlanLimit struct {
	Id          int64       `json:"id"          description:"Id"`                    // Id
	MerchantId  int64       `json:"merchantId"  description:"merchantId"`            // merchantId
	MetricId    int64       `json:"metricId"    description:"metricId"`              // metricId
	PlanId      int64       `json:"planId"      description:"plan_id"`               // plan_id
	MetricLimit int64       `json:"metricLimit" description:"plan metric limit"`     // plan metric limit
	GmtCreate   *gtime.Time `json:"gmtCreate"   description:"create time"`           // create time
	GmtModify   *gtime.Time `json:"gmtModify"   description:"update time"`           // update time
	IsDeleted   int         `json:"isDeleted"   description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	CreateTime  int64       `json:"createTime"  description:"create utc time"`       // create utc time
}
