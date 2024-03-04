package invoice

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type ListReq struct {
	g.Meta        `path:"/list" tags:"User-Invoice" method:"get,post" summary:"Invoice List"`
	UserId        int    `json:"userId" dc:"UserId Filter, Default Filter All" `
	SendEmail     string `json:"sendEmail" dc:"SendEmail Filter , Default Filter All" `
	SortField     string `json:"sortField" dc:"Filter，em. invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType      string `json:"sortType" dc:"Sort，asc|desc，Default desc" `
	DeleteInclude bool   `json:"deleteInclude" dc:"Deleted Involved，Need Admin" `
	Page          int    `json:"page"  dc:"Page, Start 0" `
	Count         int    `json:"count"  dc:"Count" dc:"Count By Page" `
}

type ListRes struct {
	Invoices []*ro.InvoiceDetailRo `json:"invoices" dc:"Invoices Details"`
}
