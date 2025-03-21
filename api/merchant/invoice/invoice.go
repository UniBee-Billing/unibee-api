package invoice

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type PdfGenerateReq struct {
	g.Meta        `path:"/pdf_generate" tags:"Invoice" method:"post" summary:"Generate Invoice PDF"`
	InvoiceId     string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
	SendUserEmail bool   `json:"sendUserEmail" d:"false" dc:"Whether sen invoice email to user or not，default false"`
}
type PdfGenerateRes struct {
}

type PdfUpdateReq struct {
	g.Meta                           `path:"/pdf_update" tags:"Invoice" method:"post" summary:"Update Invoice PDF"`
	InvoiceId                        string                 `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
	IssueCompanyName                 *string                `json:"issueCompanyName" dc:"IssueCompanyName"`
	IssueAddress                     *string                `json:"issueAddress" dc:"IssueAddress"`
	IssueVatNumber                   *string                `json:"issueVatNumber" dc:"IssueVatNumber"`
	IssueRegNumber                   *string                `json:"issueRegNumber" dc:"IssueRegNumber"`
	LocalizedCurrency                *string                `json:"localizedCurrency" dc:"LocalizedCurrency, To display localized currency amount"`
	LocalizedExchangeRate            *float64               `json:"localizedExchangeRate" dc:"LocalizedExchangeRate, exchange rate must set while LocalizedCurrency enabled"`
	LocalizedExchangeRateDescription *float64               `json:"localizedExchangeRateDescription" dc:"LocalizedExchangeRateDescription"`
	ShowDetailItem                   *bool                  `json:"showDetailItem" d:"false" dc:"ShowDetailItem, whether to display detail item information in pdf generate, unitAmount, quantity, etc."`
	SendUserEmail                    bool                   `json:"sendUserEmail" d:"false" dc:"Whether sen invoice email to user or not，default false"`
	Template                         string                 `json:"template" dc:"Template"`
	Metadata                         map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}
type PdfUpdateRes struct {
}

type SendEmailReq struct {
	g.Meta    `path:"/send_email" tags:"Invoice" method:"post" summary:"Send Invoice Email"`
	InvoiceId string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
}
type SendEmailRes struct {
}

type ReconvertCryptoAndSendReq struct {
	g.Meta    `path:"/reconvert_crypto_and_send_email" tags:"Invoice" method:"post" summary:"Admin Reconvert Crypto Data and Send Invoice Email"`
	InvoiceId string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
}
type ReconvertCryptoAndSendRes struct {
}

type DetailReq struct {
	g.Meta    `path:"/detail" tags:"Invoice" method:"get,post" summary:"Invoice Detail" dc:"Get detail of invoice"`
	InvoiceId string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
}
type DetailRes struct {
	Invoice     *detail.InvoiceDetail   `json:"invoice" dc:"Invoice Detail Object"`
	CreditNotes []*detail.InvoiceDetail `json:"creditNotes" dc:"CreditNotes Object List Link To Invoice"`
}

type ListReq struct {
	g.Meta          `path:"/list" tags:"Invoice" method:"get,post" summary:"Get Invoice List" dc:"Get invoice list"`
	FirstName       string `json:"firstName" dc:"The firstName of invoice" `
	LastName        string `json:"lastName" dc:"The lastName of invoice" `
	Currency        string `json:"currency" dc:"The currency of invoice" `
	Status          []int  `json:"status" dc:"The status of invoice, 1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled" `
	AmountStart     *int64 `json:"amountStart" dc:"The filter start amount of invoice" `
	AmountEnd       *int64 `json:"amountEnd" dc:"The filter end amount of invoice" `
	UserId          uint64 `json:"userId" dc:"The filter userid of invoice" `
	SendEmail       string `json:"sendEmail" dc:"The filter email of invoice" `
	SortField       string `json:"sortField" dc:"Filter，em. invoice_id|gmt_create|gmt_modify|period_end|total_amount，Default gmt_modify" `
	SortType        string `json:"sortType" dc:"Sort，asc|desc，Default desc" `
	DeleteInclude   bool   `json:"deleteInclude" dc:"Deleted Involved，Need Admin Permission" `
	Type            *int   `json:"type"  dc:"invoice Type, 0-payment, 1-refund" `
	Page            int    `json:"page"  dc:"Page, Start 0" `
	Count           int    `json:"count"  dc:"Count" dc:"Count By Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
	ReportTimeStart int64  `json:"reportTimeStart" dc:"ReportTimeStart" `
	ReportTimeEnd   int64  `json:"reportTimeEnd" dc:"ReportTimeEnd" `
}

type ListRes struct {
	Invoices []*detail.InvoiceDetail `json:"invoices" dc:"Invoice Detail Object List"`
	Total    int                     `json:"total" dc:"Total"`
}

type NewReq struct {
	g.Meta        `path:"/new" tags:"Invoice" method:"post" summary:"New Invoice"`
	UserId        uint64                 `json:"userId" dc:"The userId of invoice" v:"required"`
	TaxPercentage int64                  `json:"taxPercentage"  dc:"The tax percentage of invoice，1000=10%" v:"required" `
	GatewayId     uint64                 `json:"gatewayId" dc:"The gateway id of invoice" `
	Currency      string                 `json:"currency"   dc:"The currency of invoice" v:"required" `
	Name          string                 `json:"name"   dc:"The name of invoice" `
	Lines         []*NewInvoiceItemParam `json:"lines"              `
	Finish        bool                   `json:"finish" `
}

type NewInvoiceItemParam struct {
	UnitAmountExcludingTax int64  `json:"unitAmountExcludingTax"`
	Name                   string `json:"name"`
	Description            string `json:"description"`
	Quantity               int64  `json:"quantity"`
}

type NewRes struct {
	Invoice *detail.InvoiceDetail `json:"invoice" dc:"The Invoice Detail Object"`
}

type EditReq struct {
	g.Meta        `path:"/edit" tags:"Invoice" method:"post" summary:"Invoice Edit" dc:"Edit invoice of pending status"`
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
	g.Meta    `path:"/delete" tags:"Invoice" method:"post" summary:"Delete Pending Invoice" dc:"Delete invoice of pending status"`
	InvoiceId string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
}
type DeleteRes struct {
}

type FinishReq struct {
	g.Meta    `path:"/finish" tags:"Invoice" method:"post" summary:"Finish Invoice" dc:"Finish invoice, generate invoice link"`
	InvoiceId string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
	//PayMethod   int    `json:"payMethod" dc:"PayMethod,1-manual，2-auto" v:"required"`
	DaysUtilDue int `json:"daysUtilDue" dc:"Due Day Of Invoice Payment" v:"required"`
}
type FinishRes struct {
	Invoice *bean.Invoice `json:"invoice" `
}

type CancelReq struct {
	g.Meta    `path:"/cancel" tags:"Invoice" method:"post" summary:"Admin Cancel Invoice"`
	InvoiceId string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
}
type CancelRes struct {
}

type RefundReq struct {
	g.Meta       `path:"/refund" tags:"Invoice" method:"post" summary:"Create InvoiceRefund" dc:"Create payment refund for paid invoice"`
	InvoiceId    string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
	RefundNo     string `json:"refundNo" dc:"The out refund number"`
	RefundAmount int64  `json:"refundAmount" dc:"The amount of refund" v:"required"`
	Reason       string `json:"reason" dc:"The reason of refund" v:"required"`
}

type RefundRes struct {
	Refund *bean.Refund `json:"refund" dc:"Refund Object"`
}

type MarkRefundReq struct {
	g.Meta       `path:"/mark_refund" tags:"Invoice" method:"post" summary:"Mark Invoice Refund" dc:"Mark invoice as refund"`
	InvoiceId    string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
	RefundNo     string `json:"refundNo" dc:"The out refund number"`
	RefundAmount int64  `json:"refundAmount" dc:"The amount of refund" v:"required"`
	Reason       string `json:"reason" dc:"The reason of refund" v:"required"`
}

type MarkRefundRes struct {
	Refund *bean.Refund `json:"refund" dc:"Refund Object"`
}

type MarkWireTransferSuccessReq struct {
	g.Meta         `path:"/mark_wire_transfer_success" tags:"Invoice" method:"post" summary:"Mark Wire Transfer Invoice As Success" dc:"Mark wire transfer pending invoice as success"`
	InvoiceId      string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
	TransferNumber string `json:"transferNumber" dc:"The transfer number of invoice" v:"required"`
	Reason         string `json:"reason" dc:"The reason of mark action"`
}

type MarkWireTransferSuccessRes struct {
}

type MarkRefundInvoiceSuccessReq struct {
	g.Meta    `path:"/mark_refund_success" tags:"Invoice" method:"post" summary:"Mark Invoice Refund As Success" dc:"Mark refund invoice success, only support Changelly and Wire Transfer"`
	InvoiceId string `json:"invoiceId" dc:"The unique id of invoice" v:"required"`
	Reason    string `json:"reason" dc:"The reason of mark action"`
}

type MarkRefundInvoiceSuccessRes struct {
}
