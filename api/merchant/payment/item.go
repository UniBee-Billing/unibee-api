package payment

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type ItemListReq struct {
	g.Meta    `path:"/item/list" tags:"Payment" method:"get" summary:"OneTimePaymentItemList"`
	UserId    uint64 `json:"userId" dc:"Filter UserId, Default All" `
	SortField string `json:"sortField" dc:"Sort，invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType  string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page      int    `json:"page"  dc:"Page,Start 0" `
	Count     int    `json:"count" dc:"Count Of Page" `
}

type ItemListRes struct {
	PaymentItems []*bean.PaymentItem `json:"paymentItems" description:"Payment Item Object List" `
	Total        int                 `json:"total" dc:"Total"`
}
