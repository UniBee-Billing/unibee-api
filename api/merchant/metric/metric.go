package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee-api/internal/logic/metric"
)

type MerchantMetricListReq struct {
	g.Meta     `path:"/merchant_metric_list" tags:"Merchant-Metric-Controller" method:"get" summary:"Merchant Metric list"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required"`
}

type MerchantMetricListRes struct {
	MerchantMetrics []*metric.MerchantMetricVo
}

type NewMerchantMetricReq struct {
	g.Meta              `path:"/new_merchant_metric" tags:"Merchant-Metric-Controller" method:"post" summary:"New Merchant Metric"`
	MerchantId          int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	Code                string `p:"code" dc:"Code" v:"required"`
	Name                string `p:"name" dc:"Name" v:"required"`
	Description         string `p:"description" dc:"Description"`
	AggregationType     int    `p:"aggregationType" dc:"AggregationType,0-count，1-count unique, 2-latest, 3-max, 4-sum"`
	AggregationProperty string `p:"aggregationProperty" dc:"AggregationProperty, Will Needed When AggregationType != count"`
}

type NewMerchantMetricRes struct {
	MerchantMetric *metric.MerchantMetricVo
}

type EditMerchantMetricReq struct {
	g.Meta      `path:"/edit_merchant_metric" tags:"Merchant-Metric-Controller" method:"post" summary:"Edit Merchant Metric"`
	MerchantId  int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	MetricId    int64  `p:"metricId" dc:"MetricId" v:"required"`
	Name        string `p:"name" dc:"Name" v:"required"`
	Description string `p:"description" dc:"Description"`
}

type EditMerchantMetricRes struct {
	MerchantMetric *metric.MerchantMetricVo
}

type DelMerchantMetricReq struct {
	g.Meta     `path:"/delete_merchant_metric" tags:"Merchant-Metric-Controller" method:"post" summary:"Delete Merchant Metric"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required"`
	MetricId   int64 `p:"metricId" dc:"MetricId" v:"required"`
}

type DelMerchantMetricRes struct {
	MerchantMetric *metric.MerchantMetricVo
}
