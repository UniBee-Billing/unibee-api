package payment

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

type TimeLineListReq struct {
	g.Meta    `path:"/payment_timeline_list" tags:"Merchant-Payment-Timeline-Controller" method:"post" summary:"Payment TimeLine List"`
	UserId    int    `p:"userId" dc:"Filter UserId, Default All" `
	SortField string `p:"sortField" dc:"Sort，invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType  string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `p:"page"  dc:"Page,Start 0" `
	Count     int    `p:"count" dc:"Count Of Page" `
}

type TimeLineListRes struct {
	PaymentTimeLines []*entity.PaymentTimeline `json:"paymentTimeLines" description:"PaymentTimeLines" `
}
