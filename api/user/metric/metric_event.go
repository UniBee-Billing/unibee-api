package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type UserMetricReq struct {
	g.Meta    `path:"/metric" tags:"User-Metric-Event" method:"get" summary:"Query User Metric"`
	ProductId int64 `json:"productId" dc:"Id of product" dc:"default product will use if productId not specified and subscriptionId is blank"`
}

type UserMetricRes struct {
	UserMetric *detail.UserMetric `json:"userMetric" dc:"UserMetric"`
}

type UserSubscriptionMetricReq struct {
	g.Meta         `path:"/sub/metric" tags:"User-Metric-Event" method:"get" summary:"Query User Metric By Subscription"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId"`
}

type UserSubscriptionMetricRes struct {
	UserMetric *detail.UserMetric `json:"userMetric" dc:"UserMetric"`
}

type EventListReq struct {
	g.Meta          `path:"/event_list" tags:"User-Metric-Event" method:"get,post" summary:"User Metric Event List"`
	MetricIds       []int64 `json:"metricIds" dc:"Filter MetricIds, Default All" `
	SortField       string  `json:"sortField" dc:"Sort，user_id|gmt_create，Default gmt_create" `
	SortType        string  `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int     `json:"page"  dc:"Page,Start 0" `
	Count           int     `json:"count" dc:"Count OF Page" `
	CreateTimeStart int64   `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64   `json:"createTimeEnd" dc:"CreateTimeEnd" `
}

type EventListRes struct {
	Events []*detail.MerchantMetricEventDetail `json:"events" description:"User Metric Event List" `
	Total  int                                 `json:"total" dc:"Total"`
}
