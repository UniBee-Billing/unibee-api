package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type ListReq struct {
	g.Meta          `path:"/list" tags:"Metric" method:"get,post" summary:"Get Merchant Metric list"`
	SortField       string `json:"sortField" dc:"Sort，user_id|gmt_create，Default gmt_create" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page,Start 0" `
	Count           int    `json:"count" dc:"Count OF Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
}

type ListRes struct {
	MerchantMetrics []*bean.MerchantMetric `json:"merchantMetrics" dc:"MerchantMetrics"`
	Total           int                    `json:"total" dc:"Total"`
}

type NewReq struct {
	g.Meta              `path:"/new" tags:"Metric" method:"post" summary:"New Merchant Metric"`
	Code                string `json:"code" dc:"Code" v:"required"`
	Type                *int   `json:"type"                description:"1-limit_metered，2-charge_metered,3-charge_recurring"`
	MetricName          string `json:"metricName" dc:"MetricName" v:"required"`
	MetricDescription   string `json:"metricDescription" dc:"MetricDescription"`
	AggregationType     int    `json:"aggregationType" dc:"AggregationType,1-count，2-count unique, 3-latest, 4-max, 5-sum"`
	AggregationProperty string `json:"aggregationProperty" dc:"AggregationProperty, Will Needed When AggregationType != count"`
}

type NewRes struct {
	MerchantMetric *bean.MerchantMetric `json:"merchantMetric" dc:"MerchantMetric"`
}

type EditReq struct {
	g.Meta            `path:"/edit" tags:"Metric" method:"post" summary:"Edit Merchant Metric"`
	MetricId          uint64 `json:"metricId" dc:"MetricId" v:"required"`
	Type              *int   `json:"type"                description:"1-limit_metered，2-charge_metered,3-charge_recurring"`
	MetricName        string `json:"metricName" dc:"MetricName" v:"required"`
	MetricDescription string `json:"metricDescription" dc:"MetricDescription"`
}

type EditRes struct {
	MerchantMetric *bean.MerchantMetric `json:"merchantMetric" dc:"MerchantMetric"`
}

type DeleteReq struct {
	g.Meta   `path:"/delete" tags:"Metric" method:"post" summary:"Delete Merchant Metric"`
	MetricId uint64 `json:"metricId" dc:"MetricId" v:"required"`
}

type DeleteRes struct {
	MerchantMetric *bean.MerchantMetric `json:"merchantMetric" dc:"MerchantMetric"`
}

type DetailReq struct {
	g.Meta   `path:"/detail" tags:"Metric" method:"post" summary:"Merchant Metric Detail"`
	MetricId uint64 `json:"metricId" dc:"MetricId" v:"required"`
}

type DetailRes struct {
	MerchantMetric *bean.MerchantMetric `json:"merchantMetric" dc:"MerchantMetric"`
}
