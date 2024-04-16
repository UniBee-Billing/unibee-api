package invoice

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type ListReq struct {
	g.Meta    `path:"/list" tags:"User-Invoice" method:"get,post" summary:"Invoice List"`
	SortField string `json:"sortField" dc:"Filter，em. invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType  string `json:"sortType" dc:"Sort，asc|desc，Default desc" `
	Page      int    `json:"page"  dc:"Page, Start 0" `
	Count     int    `json:"count"  dc:"Count" dc:"Count By Page" `
}

type ListRes struct {
	Invoices []*detail.InvoiceDetail `json:"invoices" dc:"Invoices Details"`
}
