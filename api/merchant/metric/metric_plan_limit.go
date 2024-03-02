package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type NewPlanLimitReq struct {
	g.Meta      `path:"/plan/limit/new" tags:"Merchant-Metric-Controller" method:"post" summary:"New Merchant Metric Plan TotalLimit"`
	MetricId    int64  `json:"metricId" dc:"MetricId" v:"required"`
	PlanId      uint64 `json:"planId" dc:"PlanId" v:"required"`
	MetricLimit uint64 `json:"metricLimit" dc:"MetricLimit" v:"required"`
}

type NewPlanLimitRes struct {
	MerchantMetricPlanLimit *ro.MerchantMetricPlanLimitVo `json:"merchantMetricPlanLimit" dc:"MerchantMetricPlanLimit"`
}

type EditPlanLimitReq struct {
	g.Meta            `path:"/plan/limit/edit" tags:"Merchant-Metric-Controller" method:"post" summary:"Edit Merchant Metric Plan TotalLimit"`
	MetricPlanLimitId int64  `json:"metricPlanLimitId" dc:"MetricPlanLimitId" v:"required"`
	MetricLimit       uint64 `json:"metricLimit" dc:"MetricLimit" v:"required"`
}

type EditPlanLimitRes struct {
	MerchantMetricPlanLimit *ro.MerchantMetricPlanLimitVo `json:"merchantMetricPlanLimit" dc:"MerchantMetricPlanLimit"`
}

type DeletePlanLimitReq struct {
	g.Meta            `path:"/plan/limit/delete" tags:"Merchant-Metric-Controller" method:"post" summary:"Delete Merchant Metric Plan TotalLimit"`
	MetricPlanLimitId int64 `json:"metricPlanLimitId" dc:"MetricPlanLimitId" v:"required"`
}

type DeletePlanLimitRes struct {
	MerchantMetricPlanLimit *ro.MerchantMetricPlanLimitVo `json:"merchantMetricPlanLimit" dc:"MerchantMetricPlanLimit"`
}
