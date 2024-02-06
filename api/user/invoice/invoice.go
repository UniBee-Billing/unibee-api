package invoice

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee-api/internal/logic/gateway/ro"
)

type SubscriptionInvoiceListReq struct {
	g.Meta        `path:"/subscription_invoice_list" tags:"User-Invoice-Controller" method:"post" summary:"Invoice List"`
	MerchantId    int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	UserId        int    `p:"userId" dc:"UserId Filter, Default Filter All" `
	SendEmail     int    `p:"sendEmail" dc:"SendEmail Filter , Default Filter All" `
	SortField     string `p:"sortField" dc:"Filter，em. invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType      string `p:"sortType" dc:"Sort，asc|desc，Default desc" `
	DeleteInclude bool   `p:"deleteInclude" dc:"Deleted Involved，Need Admin" `
	Page          int    `p:"page"  dc:"Page, Start 0" `
	Count         int    `p:"count"  dc:"Count" dc:"Count By Page" `
}

type SubscriptionInvoiceListRes struct {
	Invoices []*ro.InvoiceDetailRo `p:"invoices" dc:"Invoices Details"`
}
