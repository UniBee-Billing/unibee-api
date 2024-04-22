package payment

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type TimeLineListReq struct {
	g.Meta    `path:"/payment_timeline_list" tags:"User-Payment-Timeline" method:"get" summary:"PaymentTimeLine List"`
	SortField string `json:"sortField" dc:"Sort Field，invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType  string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `json:"page"  dc:"Page, Start With 0" `
	Count     int    `json:"count" dc:"Count Of Page" `
}

type TimeLineListRes struct {
	PaymentTimelines []*bean.PaymentTimelineSimplify `json:"paymentTimeline" dc:"PaymentTimelines"`
}
