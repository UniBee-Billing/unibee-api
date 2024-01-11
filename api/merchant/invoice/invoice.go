package invoice

import "github.com/gogf/gf/v2/frame/g"

type SubscriptionInvoicePdfGenerateReq struct {
	g.Meta        `path:"/subscription_invoice_pdf_generate" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin 操作 Merchant Invoice 生成 pdf"`
	InvoiceId     string `p:"invoiceId" dc:"Invoice ID" v:"required#请输入InvoiceId"`
	SendUserEmail bool   `p:"sendUserEmail" d:"false" dc:"是否发送 Invoice Email 给到 User，默认不发送"`
}
type SubscriptionInvoicePdfGenerateRes struct {
}

type SubscriptionInvoiceSendEmailReq struct {
	g.Meta    `path:"/subscription_invoice_send_user_email" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin 操作 Merchant Invoice 发送 Email 给 User"`
	InvoiceId string `p:"invoiceId" dc:"Invoice ID" v:"required#请输入InvoiceId"`
}
type SubscriptionInvoiceSendEmailRes struct {
}
