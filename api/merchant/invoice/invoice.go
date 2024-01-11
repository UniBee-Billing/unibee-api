package invoice

import "github.com/gogf/gf/v2/frame/g"

type SubscriptionInvoicePdfGenerateReq struct {
	g.Meta    `path:"/subscription_invoice_pdf_generate" tags:"Merchant-Invoice-Controller" method:"post" summary:"Merchant Invoice 生成 pdf"`
	InvoiceId string `p:"invoiceId" dc:"Invoice ID" v:"required#请输入InvoiceId"`
}
type SubscriptionInvoicePdfGenerateRes struct {
}
