package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type ListReq struct {
	g.Meta `path:"/list" tags:"Merchant-Metric-Controller" method:"get" summary:"Merchant Metric list"`
}

type ListRes struct {
	MerchantMetrics []*ro.MerchantMetricVo `json:"merchantMetrics" dc:"MerchantMetrics"`
}

type NewReq struct {
	g.Meta              `path:"/new" tags:"Merchant-Metric-Controller" method:"post" summary:"New Merchant Metric"`
	Code                string `p:"code" dc:"Code" v:"required"`
	MetricName          string `p:"metricName" dc:"MetricName" v:"required"`
	MetricDescription   string `p:"metricDescription" dc:"MetricDescription"`
	AggregationType     int    `p:"aggregationType" dc:"AggregationType,1-count，2-count unique, 3-latest, 4-max, 5-sum"`
	AggregationProperty string `p:"aggregationProperty" dc:"AggregationProperty, Will Needed When AggregationType != count"`
}

type NewRes struct {
	MerchantMetric *ro.MerchantMetricVo `json:"merchantMetric" dc:"MerchantMetric"`
}

type EditReq struct {
	g.Meta            `path:"/edit" tags:"Merchant-Metric-Controller" method:"post" summary:"Edit Merchant Metric"`
	MetricId          int64  `p:"metricId" dc:"MetricId" v:"required"`
	MetricName        string `p:"metricName" dc:"MetricName" v:"required"`
	MetricDescription string `p:"metricDescription" dc:"MetricDescription"`
}

type EditRes struct {
	MerchantMetric *ro.MerchantMetricVo `json:"merchantMetric" dc:"MerchantMetric"`
}

type DelReq struct {
	g.Meta   `path:"/delete" tags:"Merchant-Metric-Controller" method:"post" summary:"Delete Merchant Metric"`
	MetricId int64 `p:"metricId" dc:"MetricId" v:"required"`
}

type DelRes struct {
	MerchantMetric *ro.MerchantMetricVo `json:"merchantMetric" dc:"MerchantMetric"`
}

type DetailReq struct {
	g.Meta   `path:"/detail" tags:"Merchant-Metric-Controller" method:"post" summary:"Merchant Metric Detail"`
	MetricId uint64 `p:"metricId" dc:"MetricId" v:"required"`
}

type DetailRes struct {
	MerchantMetric *ro.MerchantMetricVo `json:"merchantMetric" dc:"MerchantMetric"`
}
