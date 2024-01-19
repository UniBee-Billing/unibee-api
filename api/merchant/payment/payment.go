package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type EventListReq struct {
	g.Meta     `path:"/payment_event_list" tags:"Merchant-Payment-Event-Controller" method:"post" summary:"PaymentEvent列表"`
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	UserId     int    `p:"userId" dc:"UserId 不填查询所有" `
	SortField  string `p:"sortField" dc:"排序字段，invoice_id|gmt_create|gmt_modify|period_end|total_amount，默认 gmt_modify" `
	SortType   string `p:"sortType" dc:"排序类型，asc|desc，默认 desc" `
	Page       int    `p:"page"  dc:"分页页码,0开始" `
	Count      int    `p:"count"  dc:"订阅计划货币" dc:"每页数量" `
}

type EventListRes struct {
}
