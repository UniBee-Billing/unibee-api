package metric

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type EventListReq struct {
	g.Meta          `path:"/event_list" tags:"User-Metric-Event" method:"get,post" summary:"User Metric Event List"`
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
