package metric

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type NewEventReq struct {
	g.Meta              `path:"/event/new" tags:"Metric Event" method:"post" summary:"New Merchant Metric Event"`
	MetricCode          string      `json:"metricCode" dc:"MetricCode" v:"required"`
	ExternalEventId     string      `json:"externalEventId" dc:"ExternalEventId, __unique__" v:"required"`
	ExternalUserId      string      `json:"externalUserId" dc:"ExternalUserId， UserId, ExternalUserId, or Email provides one of three options" `
	Email               string      `json:"email" dc:"Email， UserId, ExternalUserId, or Email provides one of three options" default:"account@unibee.dev"`
	UserId              uint64      `json:"userId" dc:"UserId， UserId, ExternalUserId, or Email provides one of three options" `
	MetricProperties    *gjson.Json `json:"metricProperties" dc:"Metric property data,Json format, find property data on it when AggregationValue or AggregationUniqueId not specified"`
	AggregationValue    *uint64     `json:"aggregationValue" dc:"AggregationValue, valid when AggregationType latest, max or sum"`
	AggregationUniqueId *string     `json:"aggregationUniqueId" dc:"AggregationUniqueId, valid when AggregationType is count unique"`
	ProductId           int64       `json:"productId" dc:"Id of product" dc:"default product will use if productId not specified and subscriptionId is blank"`
}

type NewEventRes struct {
	MerchantMetricEvent *bean.MerchantMetricEvent `json:"merchantMetricEvent" dc:"MerchantMetricEvent"`
}

type DeleteEventReq struct {
	g.Meta          `path:"/event/delete" tags:"Metric Event" method:"post" summary:"Del Merchant Metric Event"`
	MetricCode      string `json:"metricCode" dc:"MetricCode" v:"required"`
	ExternalUserId  string `json:"externalUserId" dc:"ExternalUserId， UserId, ExternalUserId, or Email provides one of three options" `
	Email           string `json:"email" dc:"Email， UserId,ExternalUserId, or Email provides one of three options" `
	UserId          uint64 `json:"userId" dc:"UserId， UserId,ExternalUserId, or Email provides one of three options" `
	ExternalEventId string `json:"externalEventId" dc:"ExternalEventId" v:"required"`
}

type DeleteEventRes struct {
}

type EventListReq struct {
	g.Meta          `path:"/event_list" tags:"Metric Event" method:"get,post" summary:"Metric Event List"`
	UserIds         []int64 `json:"userIds" dc:"Filter UserIds, Default All" `
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
