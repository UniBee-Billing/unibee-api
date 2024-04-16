package invoice

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type ListReq struct {
	g.Meta      `path:"/list" tags:"User-Invoice" method:"get,post" summary:"Invoice List"`
	Currency    string `json:"currency" dc:"The currency of invoice" `
	Status      []int  `json:"status" dc:"The status of invoice, 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled" `
	AmountStart int64  `json:"amountStart" dc:"The filter start amount of invoice" `
	AmountEnd   int64  `json:"amountEnd" dc:"The filter end amount of invoice" `
	SortField   string `json:"sortField" dc:"Filter，em. invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType    string `json:"sortType" dc:"Sort，asc|desc，Default desc" `
	Page        int    `json:"page"  dc:"Page, Start 0" `
	Count       int    `json:"count"  dc:"Count" dc:"Count By Page" `
}

type ListRes struct {
	Invoices []*detail.InvoiceDetail `json:"invoices" dc:"Invoices Details"`
}
