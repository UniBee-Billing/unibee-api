package invoice

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee-api/internal/logic/gateway/ro"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

type SubscriptionInvoicePdfGenerateReq struct {
	g.Meta        `path:"/subscription_invoice_pdf_generate" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Generate Merchant Invoice pdf"`
	InvoiceId     string `p:"invoiceId" dc:"Invoice ID" v:"required"`
	SendUserEmail bool   `p:"sendUserEmail" d:"false" dc:"Whether Send Invoice Email To User Or Not，Default Not Send"`
}
type SubscriptionInvoicePdfGenerateRes struct {
}

type SubscriptionInvoiceSendEmailReq struct {
	g.Meta    `path:"/subscription_invoice_send_user_email" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Send Merchant Invoice Email to User"`
	InvoiceId string `p:"invoiceId" dc:"Invoice ID" v:"required"`
}
type SubscriptionInvoiceSendEmailRes struct {
}

type SubscriptionInvoiceDetailReq struct {
	g.Meta    `path:"/subscription_invoice_detail" tags:"Merchant-Invoice-Controller" method:"post" summary:"Invoice Detail"`
	InvoiceId string `p:"invoiceId" dc:"Invoice ID" v:"required"`
}
type SubscriptionInvoiceDetailRes struct {
	Invoice *ro.InvoiceDetailRo `p:"invoice" dc:"invoice Detail"`
}

type SubscriptionInvoiceListReq struct {
	g.Meta        `path:"/subscription_invoice_list" tags:"Merchant-Invoice-Controller" method:"post" summary:"Invoice List"`
	FirstName     string `p:"firstName" dc:"FirstName" `
	LastName      string `p:"lastName" dc:"LastName" `
	Currency      string `p:"currency" dc:"Currency" `
	Status        []int  `p:"status" dc:"Status, 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled" `
	AmountStart   int64  `p:"amountStart" dc:"AmountStart" `
	AmountEnd     int64  `p:"amountEnd" dc:"AmountEnd" `
	UserId        int    `p:"userId" dc:"UserId Filter, Default Filter All" `
	SendEmail     string `p:"sendEmail" dc:"SendEmail Filter , Default Filter All" `
	SortField     string `p:"sortField" dc:"Filter，em. invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType      string `p:"sortType" dc:"Sort，asc|desc，Default desc" `
	DeleteInclude bool   `p:"deleteInclude" dc:"Deleted Involved，Need Admin" `
	Page          int    `p:"page"  dc:"Page, Start 0" `
	Count         int    `p:"count"  dc:"Count" dc:"Count By Page" `
}

type SubscriptionInvoiceListRes struct {
	Invoices []*ro.InvoiceDetailRo `p:"invoices" dc:"invoice Detail List"`
}

type NewInvoiceCreateReq struct {
	g.Meta    `path:"/new_invoice_create" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Create New Invoice"`
	UserId    int64                  `p:"userId" dc:"UserId" v:"required"`
	TaxScale  int64                  `p:"taxScale"  dc:"TaxScale，1000 represent 10%" v:"required" `
	GatewayId int64                  `p:"gatewayId" dc:"Gateway Id"   v:"required" `
	Currency  string                 `p:"currency"   dc:"Currency" v:"required" `
	Name      string                 `p:"name"   dc:"Name" `
	Lines     []*NewInvoiceItemParam `p:"lines"              `
	Finish    bool                   `p:"finish" `
}

type NewInvoiceItemParam struct {
	UnitAmountExcludingTax int64  `json:"unitAmountExcludingTax"`
	Description            string `json:"description"`
	Quantity               int64  `json:"quantity"`
}

type NewInvoiceCreateRes struct {
	Invoice *ro.InvoiceDetailRo `json:"invoice" `
}

type NewInvoiceEditReq struct {
	g.Meta    `path:"/new_invoice_edit" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Edit Invoice"`
	InvoiceId string                 `p:"invoiceId" dc:"InvoiceId" v:"required|length:4,30#请输入InvoiceId"`
	TaxScale  int64                  `p:"taxScale"  dc:"TaxScale，1000 represent 10%"`
	GatewayId int64                  `p:"gatewayId" dc:"Gateway Id" `
	Currency  string                 `p:"currency"   dc:"Currency" `
	Name      string                 `p:"name"   dc:"Name" `
	Lines     []*NewInvoiceItemParam `p:"lines"              `
	Finish    bool                   `p:"finish" `
}
type NewInvoiceEditRes struct {
	Invoice *ro.InvoiceDetailRo `json:"invoice" `
}

type DeletePendingInvoiceReq struct {
	g.Meta    `path:"/new_invoice_delete" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Delete Invoice Of Pending Status"`
	InvoiceId string `p:"invoiceId" dc:"InvoiceId" v:"required"`
}
type DeletePendingInvoiceRes struct {
}

type ProcessInvoiceForPayReq struct {
	g.Meta      `path:"/finish_new_invoice" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Finish Invoice，Generate Pay Link"`
	InvoiceId   string `p:"invoiceId" dc:"InvoiceId" v:"required"`
	PayMethod   int    `p:"payMethod" dc:"PayMethod,1-manual，2-auto" v:"required"`
	DaysUtilDue int    `p:"daysUtilDue" dc:"DaysUtilDue,Due Day Of Pay" v:"required"`
}
type ProcessInvoiceForPayRes struct {
	Invoice *entity.Invoice `json:"invoice" `
}

type CancelProcessingInvoiceReq struct {
	g.Meta    `path:"/cancel_processing_invoice" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Cancel Invoice Of Processing Status"`
	InvoiceId string `p:"invoiceId" dc:"InvoiceId" v:"required`
}
type CancelProcessingInvoiceRes struct {
}

type NewInvoiceRefundReq struct {
	g.Meta       `path:"/new_invoice_refund" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin Create Refund From Invoice"`
	InvoiceId    string `p:"invoiceId" dc:"InvoiceId" v:"required"`
	RefundNo     string `p:"refundNo" dc:"RefundNo" v:"required"`
	RefundAmount int64  `p:"refundAmount" dc:"Refund Amount" v:"required"`
	Reason       string `p:"reason" dc:"Refund Reason" v:"required"`
}

type NewInvoiceRefundRes struct {
	Refund *entity.Refund `json:"refund" `
}
