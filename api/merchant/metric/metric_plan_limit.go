package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

type SetMerchantMetricPlanLimitReq struct {
	g.Meta     `path:"/edit_merchant_metric_plan_limit" tags:"Merchant-Metric-Controller" method:"post" summary:"Set Merchant Metric Plan Limit"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required"`
	MetricId   int64 `p:"metricId" dc:"MetricId" v:"required"`
	PlanId     int64 `p:"planId" dc:"PlanId" v:"required"`
	Limit      int64 `p:"limit" dc:"Limit" v:"required"`
}

type SetMerchantMetricPlanLimitRes struct {
	MerchantMetric *entity.MerchantMetric
}
