package payment

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type TimeLineListReq struct {
	g.Meta    `path:"/timeline/list" tags:"Payment" method:"get" summary:"PaymentTimeLineList"`
	UserId    int64  `json:"userId" dc:"Filter UserId, Default All" `
	SortField string `json:"sortField" dc:"Sort，invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType  string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `json:"page"  dc:"Page,Start 0" `
	Count     int    `json:"count" dc:"Count Of Page" `
}

type TimeLineListRes struct {
	PaymentTimeLines []*bean.PaymentTimelineSimplify `json:"paymentTimeLines" description:"Payment TimeLine Object List" `
}
