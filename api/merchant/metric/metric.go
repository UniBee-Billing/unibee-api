package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type MerchantMetricListReq struct {
	g.Meta `path:"/merchant_metric_list" tags:"Merchant-Metric-Controller" method:"get" summary:"Merchant Metric list"`
}

type MerchantMetricListRes struct {
	MerchantMetrics []*ro.MerchantMetricVo `json:"merchantMetrics" dc:"MerchantMetrics"`
}

type NewMerchantMetricReq struct {
	g.Meta              `path:"/new_merchant_metric" tags:"Merchant-Metric-Controller" method:"post" summary:"New Merchant Metric"`
	Code                string `p:"code" dc:"Code" v:"required"`
	MetricName          string `p:"metricName" dc:"MetricName" v:"required"`
	MetricDescription   string `p:"metricDescription" dc:"MetricDescription"`
	AggregationType     int    `p:"aggregationType" dc:"AggregationType,1-count，2-count unique, 3-latest, 4-max, 5-sum"`
	AggregationProperty string `p:"aggregationProperty" dc:"AggregationProperty, Will Needed When AggregationType != count"`
}

type NewMerchantMetricRes struct {
	MerchantMetric *ro.MerchantMetricVo `json:"merchantMetric" dc:"MerchantMetric"`
}

type EditMerchantMetricReq struct {
	g.Meta            `path:"/edit_merchant_metric" tags:"Merchant-Metric-Controller" method:"post" summary:"Edit Merchant Metric"`
	MetricId          int64  `p:"metricId" dc:"MetricId" v:"required"`
	MetricName        string `p:"metricName" dc:"MetricName" v:"required"`
	MetricDescription string `p:"metricDescription" dc:"MetricDescription"`
}

type EditMerchantMetricRes struct {
	MerchantMetric *ro.MerchantMetricVo `json:"merchantMetric" dc:"MerchantMetric"`
}

type DelMerchantMetricReq struct {
	g.Meta   `path:"/delete_merchant_metric" tags:"Merchant-Metric-Controller" method:"post" summary:"Delete Merchant Metric"`
	MetricId int64 `p:"metricId" dc:"MetricId" v:"required"`
}

type DelMerchantMetricRes struct {
	MerchantMetric *ro.MerchantMetricVo `json:"merchantMetric" dc:"MerchantMetric"`
}

type MerchantMetricDetailReq struct {
	g.Meta   `path:"/merchant_metric_detail" tags:"Merchant-Metric-Controller" method:"post" summary:"Merchant Metric Detail"`
	MetricId uint64 `p:"metricId" dc:"MetricId" v:"required"`
}

type MerchantMetricDetailRes struct {
	MerchantMetric *ro.MerchantMetricVo `json:"merchantMetric" dc:"MerchantMetric"`
}
