package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type NewPlanLimitReq struct {
	g.Meta      `path:"/plan/limit/new" tags:"Merchant-Metric-Controller" method:"post" summary:"New Merchant Metric Plan TotalLimit"`
	MetricId    int64  `p:"metricId" dc:"MetricId" v:"required"`
	PlanId      uint64 `p:"planId" dc:"PlanId" v:"required"`
	MetricLimit uint64 `p:"metricLimit" dc:"MetricLimit" v:"required"`
}

type NewPlanLimitRes struct {
	MerchantMetricPlanLimit *ro.MerchantMetricPlanLimitVo `json:"merchantMetricPlanLimit" dc:"MerchantMetricPlanLimit"`
}

type EditPlanLimitReq struct {
	g.Meta            `path:"/plan/limit/edit" tags:"Merchant-Metric-Controller" method:"post" summary:"Edit Merchant Metric Plan TotalLimit"`
	MetricPlanLimitId int64  `p:"metricPlanLimitId" dc:"MetricPlanLimitId" v:"required"`
	MetricLimit       uint64 `p:"metricLimit" dc:"MetricLimit" v:"required"`
}

type EditMerchantMetricPlanLimitRes struct {
	MerchantMetricPlanLimit *ro.MerchantMetricPlanLimitVo `json:"merchantMetricPlanLimit" dc:"MerchantMetricPlanLimit"`
}

type DeletePlanLimitReq struct {
	g.Meta            `path:"/plan/limit/delete" tags:"Merchant-Metric-Controller" method:"post" summary:"Delete Merchant Metric Plan TotalLimit"`
	MetricPlanLimitId int64 `p:"metricPlanLimitId" dc:"MetricPlanLimitId" v:"required"`
}

type DeletePlanLimitRes struct {
	MerchantMetricPlanLimit *ro.MerchantMetricPlanLimitVo `json:"merchantMetricPlanLimit" dc:"MerchantMetricPlanLimit"`
}
