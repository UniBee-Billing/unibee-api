package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee-api/internal/logic/gateway/ro"
)

type NewMerchantMetricPlanLimitReq struct {
	g.Meta      `path:"/new_merchant_metric_plan_limit" tags:"Merchant-Metric-Controller" method:"post" summary:"New Merchant Metric Plan TotalLimit"`
	MerchantId  uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
	MetricId    int64  `p:"metricId" dc:"MetricId" v:"required"`
	PlanId      int64  `p:"planId" dc:"PlanId" v:"required"`
	MetricLimit uint64 `p:"metricLimit" dc:"MetricLimit" v:"required"`
}

type NewMerchantMetricPlanLimitRes struct {
	MerchantMetric *ro.MerchantMetricPlanLimitVo
}

type EditMerchantMetricPlanLimitReq struct {
	g.Meta            `path:"/edit_merchant_metric_plan_limit" tags:"Merchant-Metric-Controller" method:"post" summary:"Edit Merchant Metric Plan TotalLimit"`
	MerchantId        uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
	MetricPlanLimitId int64  `p:"metricPlanLimitId" dc:"MetricPlanLimitId" v:"required"`
	MetricLimit       uint64 `p:"metricLimit" dc:"MetricLimit" v:"required"`
}

type EditMerchantMetricPlanLimitRes struct {
	MerchantMetric *ro.MerchantMetricPlanLimitVo
}

type DelMerchantMetricPlanLimitReq struct {
	g.Meta            `path:"/delete_merchant_metric_plan_limit" tags:"Merchant-Metric-Controller" method:"post" summary:"Delete Merchant Metric Plan TotalLimit"`
	MerchantId        uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
	MetricPlanLimitId int64  `p:"metricPlanLimitId" dc:"MetricPlanLimitId" v:"required"`
}

type DelMerchantMetricPlanLimitRes struct {
	MerchantMetric *ro.MerchantMetricPlanLimitVo
}
