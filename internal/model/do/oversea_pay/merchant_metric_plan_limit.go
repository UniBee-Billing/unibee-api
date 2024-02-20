// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantMetricPlanLimit is the golang structure of table merchant_metric_plan_limit for DAO operations like Where/Data.
type MerchantMetricPlanLimit struct {
	g.Meta      `orm:"table:merchant_metric_plan_limit, do:true"`
	Id          interface{} // Id
	MerchantId  interface{} // merchantId
	MetricId    interface{} // metricId
	PlanId      interface{} // plan_id
	MetricLimit interface{} // plan metric limit
	GmtCreate   *gtime.Time // create time
	GmtModify   *gtime.Time // update time
	IsDeleted   interface{} // 0-UnDeleted，1-Deleted
	CreateTime  interface{} // create utc time
}
