package payment

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type TimeLineListReq struct {
	g.Meta          `path:"/timeline/list" tags:"Payment" method:"get" summary:"PaymentTimeLineList"`
	UserId          uint64 `json:"userId" dc:"Filter UserId, Default All" `
	SortField       string `json:"sortField" dc:"Sort，invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page,Start 0" `
	Count           int    `json:"count" dc:"Count Of Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
}

type TimeLineListRes struct {
	PaymentTimeLines []*detail.PaymentTimelineDetail `json:"paymentTimeLines" description:"Payment TimeLine Object List" `
	Total            int                             `json:"total" dc:"Total"`
}
