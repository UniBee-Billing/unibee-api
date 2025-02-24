package metric

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type NewEventReq struct {
	g.Meta           `path:"/event/new" tags:"User Metric" method:"post" summary:"New Merchant Metric Event"`
	MetricCode       string      `json:"metricCode" dc:"MetricCode" v:"required"`
	ExternalUserId   string      `json:"externalUserId" dc:"ExternalUserId" v:"required"`
	ExternalEventId  string      `json:"externalEventId" dc:"ExternalEventId, __unique__" v:"required"`
	MetricProperties *gjson.Json `json:"metricProperties" dc:"MetricProperties"`
	ProductId        int64       `json:"productId" dc:"Id of product" dc:"default product will use if productId not specified and subscriptionId is blank"`
}

type NewEventRes struct {
	MerchantMetricEvent *bean.MerchantMetricEvent `json:"merchantMetricEvent" dc:"MerchantMetricEvent"`
}

type DeleteEventReq struct {
	g.Meta          `path:"/event/delete" tags:"User Metric" method:"post" summary:"Del Merchant Metric Event"`
	MetricCode      string `json:"metricCode" dc:"MetricCode" v:"required"`
	ExternalUserId  string `json:"externalUserId" dc:"ExternalUserId" v:"required"`
	ExternalEventId string `json:"externalEventId" dc:"ExternalEventId" v:"required"`
}

type DeleteEventRes struct {
}

type EventListReq struct {
	g.Meta          `path:"/event_list" tags:"User Metric" method:"get,post" summary:"User Metric Event List"`
	UserId          int64  `json:"userId" dc:"Filter UserId" `
	SortField       string `json:"sortField" dc:"Sort，user_id|gmt_create，Default gmt_create" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page,Start 0" `
	Count           int    `json:"count" dc:"Count OF Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
}

type EventListRes struct {
	Events []*detail.MerchantMetricEventDetail `json:"events" description:"User Metric Event List" `
	Total  int                                 `json:"total" dc:"Total"`
}
