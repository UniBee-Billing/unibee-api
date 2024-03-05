package payment

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "unibee/internal/model/entity/oversea_pay"
)

type TimeLineListReq struct {
	g.Meta    `path:"/payment_timeline_list" tags:"Merchant-Payment-Timeline" method:"get,post" summary:"Payment TimeLine List"`
	UserId    int    `json:"userId" dc:"Filter UserId, Default All" `
	SortField string `json:"sortField" dc:"Sort，invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType  string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `json:"page"  dc:"Page,Start 0" `
	Count     int    `json:"count" dc:"Count Of Page" `
}

type TimeLineListRes struct {
	PaymentTimeLines []*entity.PaymentTimeline `json:"paymentTimeLines" description:"PaymentTimeLines" `
}
