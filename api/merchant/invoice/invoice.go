package invoice

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
	entity "unibee/internal/model/entity/oversea_pay"
)

type PdfGenerateReq struct {
	g.Meta        `path:"/pdf_generate" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Generate Merchant Invoice pdf"`
	InvoiceId     string `json:"invoiceId" dc:"Invoice ID" v:"required"`
	SendUserEmail bool   `json:"sendUserEmail" d:"false" dc:"Whether Send Invoice Email To User Or Not，Default Not Send"`
}
type PdfGenerateRes struct {
}

type SendEmailReq struct {
	g.Meta    `path:"/send_email" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Send Merchant Invoice Email to User"`
	InvoiceId string `json:"invoiceId" dc:"Invoice ID" v:"required"`
}
type SendEmailRes struct {
}

type DetailReq struct {
	g.Meta    `path:"/detail" tags:"Merchant-Invoice-Controller" method:"post" summary:"Invoice Detail"`
	InvoiceId string `json:"invoiceId" dc:"Invoice ID" v:"required"`
}
type DetailRes struct {
	Invoice *ro.InvoiceDetailRo `json:"invoice" dc:"invoice Detail"`
}

type ListReq struct {
	g.Meta        `path:"/list" tags:"Merchant-Invoice-Controller" method:"post" summary:"Invoice List"`
	FirstName     string `json:"firstName" dc:"FirstName" `
	LastName      string `json:"lastName" dc:"LastName" `
	Currency      string `json:"currency" dc:"Currency" `
	Status        []int  `json:"status" dc:"Status, 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled" `
	AmountStart   int64  `json:"amountStart" dc:"AmountStart" `
	AmountEnd     int64  `json:"amountEnd" dc:"AmountEnd" `
	UserId        int    `json:"userId" dc:"UserId Filter, Default Filter All" `
	SendEmail     string `json:"sendEmail" dc:"SendEmail Filter , Default Filter All" `
	SortField     string `json:"sortField" dc:"Filter，em. invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType      string `json:"sortType" dc:"Sort，asc|desc，Default desc" `
	DeleteInclude bool   `json:"deleteInclude" dc:"Deleted Involved，Need Admin" `
	Page          int    `json:"page"  dc:"Page, Start 0" `
	Count         int    `json:"count"  dc:"Count" dc:"Count By Page" `
}

type ListRes struct {
	Invoices []*ro.InvoiceDetailRo `json:"invoices" dc:"invoice Detail List"`
}

type NewReq struct {
	g.Meta    `path:"/new" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Create New Invoice"`
	UserId    int64                  `json:"userId" dc:"UserId" v:"required"`
	TaxScale  int64                  `json:"taxScale"  dc:"TaxScale，1000 represent 10%" v:"required" `
	GatewayId uint64                 `json:"gatewayId" dc:"Gateway Id"   v:"required" `
	Currency  string                 `json:"currency"   dc:"Currency" v:"required" `
	Name      string                 `json:"name"   dc:"Name" `
	Lines     []*NewInvoiceItemParam `json:"lines"              `
	Finish    bool                   `json:"finish" `
}

type NewInvoiceItemParam struct {
	UnitAmountExcludingTax int64  `json:"unitAmountExcludingTax"`
	Description            string `json:"description"`
	Quantity               int64  `json:"quantity"`
}

type NewRes struct {
	Invoice *ro.InvoiceDetailRo `json:"invoice" `
}

type EditReq struct {
	g.Meta    `path:"/edit" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Edit Invoice"`
	InvoiceId string                 `json:"invoiceId" dc:"InvoiceId" v:"required|length:4,30#请输入InvoiceId"`
	TaxScale  int64                  `json:"taxScale"  dc:"TaxScale，1000 represent 10%"`
	GatewayId uint64                 `json:"gatewayId" dc:"Gateway Id" `
	Currency  string                 `json:"currency"   dc:"Currency" `
	Name      string                 `json:"name"   dc:"Name" `
	Lines     []*NewInvoiceItemParam `json:"lines"              `
	Finish    bool                   `json:"finish" `
}
type EditRes struct {
	Invoice *ro.InvoiceDetailRo `json:"invoice" `
}

type DeleteReq struct {
	g.Meta    `path:"/delete" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Delete Invoice Of Pending Status"`
	InvoiceId string `json:"invoiceId" dc:"InvoiceId" v:"required"`
}
type DeleteRes struct {
}

type FinishReq struct {
	g.Meta      `path:"/finish" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Finish Invoice，Generate Pay Link"`
	InvoiceId   string `json:"invoiceId" dc:"InvoiceId" v:"required"`
	PayMethod   int    `json:"payMethod" dc:"PayMethod,1-manual，2-auto" v:"required"`
	DaysUtilDue int    `json:"daysUtilDue" dc:"DaysUtilDue,Due Day Of Pay" v:"required"`
}
type FinishRes struct {
	Invoice *entity.Invoice `json:"invoice" `
}

type CancelReq struct {
	g.Meta    `path:"/cancel" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Cancel Invoice Of Processing Status"`
	InvoiceId string `json:"invoiceId" dc:"InvoiceId" v:"required`
}
type CancelRes struct {
}

type RefundReq struct {
	g.Meta       `path:"/refund" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Create Refund From Invoice"`
	InvoiceId    string `json:"invoiceId" dc:"InvoiceId" v:"required"`
	RefundNo     string `json:"refundNo" dc:"RefundNo" v:"required"`
	RefundAmount int64  `json:"refundAmount" dc:"Refund Amount" v:"required"`
	Reason       string `json:"reason" dc:"Refund Reason" v:"required"`
}

type RefundRes struct {
	Refund *entity.Refund `json:"refund" `
}
