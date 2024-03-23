package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type NewPlanLimitReq struct {
	g.Meta      `path:"/plan/limit/new" tags:"Metric" method:"post" summary:"New Merchant Metric Plan TotalLimit"`
	MetricId    uint64 `json:"metricId" dc:"MetricId" v:"required"`
	PlanId      uint64 `json:"planId" dc:"PlanId" v:"required"`
	MetricLimit uint64 `json:"metricLimit" dc:"MetricLimit" v:"required"`
}

type NewPlanLimitRes struct {
	MerchantMetricPlanLimit *bean.MerchantMetricPlanLimit `json:"merchantMetricPlanLimit" dc:"MerchantMetricPlanLimit"`
}

type EditPlanLimitReq struct {
	g.Meta            `path:"/plan/limit/edit" tags:"Metric" method:"post" summary:"Edit Merchant Metric Plan TotalLimit"`
	MetricPlanLimitId int64  `json:"metricPlanLimitId" dc:"MetricPlanLimitId" v:"required"`
	MetricLimit       uint64 `json:"metricLimit" dc:"MetricLimit" v:"required"`
}

type EditPlanLimitRes struct {
	MerchantMetricPlanLimit *bean.MerchantMetricPlanLimit `json:"merchantMetricPlanLimit" dc:"MerchantMetricPlanLimit"`
}

type DeletePlanLimitReq struct {
	g.Meta            `path:"/plan/limit/delete" tags:"Metric" method:"post" summary:"Delete Merchant Metric Plan TotalLimit"`
	MetricPlanLimitId uint64 `json:"metricPlanLimitId" dc:"MetricPlanLimitId" v:"required"`
}

type DeletePlanLimitRes struct {
	MerchantMetricPlanLimit *bean.MerchantMetricPlanLimit `json:"merchantMetricPlanLimit" dc:"MerchantMetricPlanLimit"`
}
