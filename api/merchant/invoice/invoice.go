package invoice

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type PdfGenerateReq struct {
	g.Meta        `path:"/pdf_generate" tags:"Invoice" method:"post" summary:"GenerateInvoicePDF"`
	InvoiceId     string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
	SendUserEmail bool   `json:"sendUserEmail" d:"false" dc:"Whether sen invoice email to user or not，default false"`
}
type PdfGenerateRes struct {
}

type SendEmailReq struct {
	g.Meta    `path:"/send_email" tags:"Invoice" method:"post" summary:"SendInvoiceEmail"`
	InvoiceId string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
}
type SendEmailRes struct {
}

type ReconvertCryptoAndSendReq struct {
	g.Meta    `path:"/reconvert_crypto_and_send_email" tags:"Invoice" method:"post" summary:"Admin Reconvert Crypto Data and Send Invoice Email to User"`
	InvoiceId string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
}
type ReconvertCryptoAndSendRes struct {
}

type DetailReq struct {
	g.Meta    `path:"/detail" tags:"Invoice" method:"get,post" summary:"InvoiceDetail" dc:"Get detail of invoice"`
	InvoiceId string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
}
type DetailRes struct {
	Invoice *detail.InvoiceDetail `json:"invoice" dc:"Invoice Detail Object"`
}

type ListReq struct {
	g.Meta        `path:"/list" tags:"Invoice" method:"get,post" summary:"InvoiceList" dc:"Get invoice list"`
	FirstName     string `json:"firstName" dc:"The firstName of invoice" `
	LastName      string `json:"lastName" dc:"The lastName of invoice" `
	Currency      string `json:"currency" dc:"The currency of invoice" `
	Status        []int  `json:"status" dc:"The status of invoice, 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled" `
	AmountStart   int64  `json:"amountStart" dc:"The filter start amount of invoice" `
	AmountEnd     int64  `json:"amountEnd" dc:"The filter end amount of invoice" `
	UserId        uint64 `json:"userId" dc:"The filter userid of invoice" `
	SendEmail     string `json:"sendEmail" dc:"The filter email of invoice" `
	SortField     string `json:"sortField" dc:"Filter，em. invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType      string `json:"sortType" dc:"Sort，asc|desc，Default desc" `
	DeleteInclude bool   `json:"deleteInclude" dc:"Deleted Involved，Need Admin Permission" `
	Page          int    `json:"page"  dc:"Page, Start 0" `
	Count         int    `json:"count"  dc:"Count" dc:"Count By Page" `
}

type ListRes struct {
	Invoices []*detail.InvoiceDetail `json:"invoices" dc:"Invoice Detail Object List"`
	Total    int                     `json:"total" dc:"Total"`
}

type NewReq struct {
	g.Meta        `path:"/new" tags:"Invoice" method:"post" summary:"NewInvoice"`
	UserId        uint64                 `json:"userId" dc:"The userId of invoice" v:"required"`
	TaxPercentage int64                  `json:"taxPercentage"  dc:"The tax percentage of invoice，1000=10%" v:"required" `
	GatewayId     uint64                 `json:"gatewayId" dc:"The gateway id of invoice"   v:"required" `
	Currency      string                 `json:"currency"   dc:"The currency of invoice" v:"required" `
	Name          string                 `json:"name"   dc:"The name of invoice" `
	Lines         []*NewInvoiceItemParam `json:"lines"              `
	Finish        bool                   `json:"finish" `
}

type NewInvoiceItemParam struct {
	UnitAmountExcludingTax int64  `json:"unitAmountExcludingTax"`
	Description            string `json:"description"`
	Quantity               int64  `json:"quantity"`
}

type NewRes struct {
	Invoice *detail.InvoiceDetail `json:"invoice" dc:"The Invoice Detail Object"`
}

type EditReq struct {
	g.Meta        `path:"/edit" tags:"Invoice" method:"post" summary:"InvoiceEdit" dc:"Edit invoice of pending status"`
	InvoiceId     string                 `json:"invoiceId" dc:"The unique id of invoice" v:"required#Invalid InvoiceId"`
	TaxPercentage int64                  `json:"taxPercentage"  dc:"The tax percentage of invoice，1000=10%"`
	GatewayId     uint64                 `json:"gatewayId" dc:"The gateway id of invoice" `
	Currency      string                 `json:"currency"   dc:"The currency of invoice" `
	Name          string                 `json:"name"   dc:"The name of invoice" `
	Lines         []*NewInvoiceItemParam `json:"lines"              `
	Finish        bool                   `json:"finish" `
}
type EditRes struct {
	Invoice *detail.InvoiceDetail `json:"invoice" dc:"The Invoice Detail Object"`
}

type DeleteReq struct {
	g.Meta    `path:"/delete" tags:"Invoice" method:"post" summary:"DeletePendingInvoice" dc:"Delete invoice of pending status"`
	InvoiceId string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
}
type DeleteRes struct {
}

type FinishReq struct {
	g.Meta    `path:"/finish" tags:"Invoice" method:"post" summary:"FinishInvoice" dc:"Finish invoice, generate invoice link"`
	InvoiceId string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
	//PayMethod   int    `json:"payMethod" dc:"PayMethod,1-manual，2-auto" v:"required"`
	DaysUtilDue int `json:"daysUtilDue" dc:"Due Day Of Invoice Payment" v:"required"`
}
type FinishRes struct {
	Invoice *bean.InvoiceSimplify `json:"invoice" `
}

type CancelReq struct {
	g.Meta    `path:"/cancel" tags:"Invoice" method:"post" summary:"Admin Cancel Invoice Of Processing Status"`
	InvoiceId string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
}
type CancelRes struct {
}

type RefundReq struct {
	g.Meta       `path:"/refund" tags:"Invoice" method:"post" summary:"CreateInvoiceRefund" dc:"Create payment refund for paid invoice"`
	InvoiceId    string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
	RefundNo     string `json:"refundNo" dc:"The out refund number"`
	RefundAmount int64  `json:"refundAmount" dc:"The amount of refund" v:"required"`
	Reason       string `json:"reason" dc:"The reason of refund" v:"required"`
}

type RefundRes struct {
	Refund *bean.RefundSimplify `json:"refund" dc:"Refund Object"`
}

type MarkRefundReq struct {
	g.Meta       `path:"/mark_refund" tags:"Invoice" method:"post" summary:"MarkInvoiceRefund" dc:"Mark invoice as refund"`
	InvoiceId    string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
	RefundNo     string `json:"refundNo" dc:"The out refund number"`
	RefundAmount int64  `json:"refundAmount" dc:"The amount of refund" v:"required"`
	Reason       string `json:"reason" dc:"The reason of refund" v:"required"`
}

type MarkRefundRes struct {
	Refund *bean.RefundSimplify `json:"refund" dc:"Refund Object"`
}

type MarkWireTransferSuccessReq struct {
	g.Meta         `path:"/mark_wire_transfer_success" tags:"Invoice" method:"post" summary:"MarkWireTransferInvoiceSuccess" dc:"Mark wire transfer pending invoice as success"`
	InvoiceId      string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
	TransferNumber string `json:"transferNumber" dc:"The transfer number of invoice" v:"required"`
	Reason         string `json:"reason" dc:"The reason of mark action"`
}

type MarkWireTransferSuccessRes struct {
}
