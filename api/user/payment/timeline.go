package payment

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type TimeLineListReq struct {
	g.Meta     `path:"/payment_timeline_list" tags:"User-Payment-Timeline-Controller" method:"post" summary:"PaymentTimeLine列表"`
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	UserId     int    `p:"userId" dc:"Filter UserId, Default All " `
	SortField  string `p:"sortField" dc:"排序字段，invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType   string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page       int    `p:"page"  dc:"Page, Start WIth 0" `
	Count      int    `p:"count" dc:"Count Of Page" `
}

type TimeLineListRes struct {
	PaymentTimelines []*entity.PaymentTimeline `p:"paymentTimeline" dc:"paymentTimelines明细"`
}
