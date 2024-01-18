package invoice

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

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

type SubscriptionInvoiceListReq struct {
	g.Meta        `path:"/subscription_invoice_list" tags:"Merchant-Invoice-Controller" method:"post" summary:"Invoice列表"`
	MerchantId    int64  `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	UserId        int    `p:"userId" dc:"UserId 不填查询所有" `
	SendEmail     int    `p:"sendEmail" dc:"SendEmail 不填查询所有" `
	SortField     string `p:"sortField" dc:"排序字段，invoice_id|gmt_create|gmt_modify|period_end|total_amount，默认 gmt_modify" `
	SortType      string `p:"sortType" dc:"排序类型，asc|desc，默认 desc" `
	DeleteInclude bool   `p:"deleteInclude" dc:"是否包含删除，查看已删除发票需要超级管理员权限" `
	Page          int    `p:"page"  dc:"分页页码,0开始" `
	Count         int    `p:"count"  dc:"订阅计划货币" dc:"每页数量" `
}

type SubscriptionInvoiceListRes struct {
	Invoices []*ro.InvoiceDetailRo `p:"invoices" dc:"invoices明细"`
}

type NewInvoiceCreateReq struct {
	g.Meta        `path:"/new_invoice_create" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin 创建新发票"`
	MerchantId    int64                `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入MerchantId"`
	UserId        int64                `p:"userId" dc:"UserId" v:"required#请输入userId"`
	TaxPercentage float64              `p:"taxPercentage"  dc:"Tax税率，万分位，1000 表示 10%" v:"required#请输入TaxPercentage" `
	ChannelId     int64                `p:"channelId" dc:"支付通道 ID"   v:"required#请输入 ChannelId" `
	Currency      string               `p:"currency"   dc:"订阅计划货币" v:"required#请输入订阅计划货币" ` // 货币
	Lines         []*ro.NewInvoiceItem `p:"lines"              `
}
type NewInvoiceCreateRes struct {
	Invoice *entity.Invoice `json:"invoice" `
}

type NewInvoiceEditReq struct {
	g.Meta        `path:"/new_invoice_edit" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin 修改新发票"`
	InvoiceId     string               `p:"invoiceId" dc:"invoiceId" v:"required|length:4,30#请输入InvoiceId"`
	TaxPercentage int64                `p:"taxPercentage"  dc:"Tax税率，万分位，1000 表示 10%"`
	ChannelId     int64                `p:"channelId" dc:"支付通道 ID" `
	Currency      string               `p:"currency"   dc:"订阅计划货币" `
	Lines         []*ro.NewInvoiceItem `p:"lines"              `
}
type NewInvoiceEditRes struct {
	Invoice *entity.Invoice `json:"invoice" `
}

type DeletePendingInvoiceReq struct {
	g.Meta    `path:"/new_invoice_delete" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin 删除Pending 状态的发票"`
	InvoiceId string `p:"invoiceId" dc:"invoiceId" v:"required|length:4,30#请输入InvoiceId"`
}
type DeletePendingInvoiceRes struct {
}

type ProcessInvoiceForPayReq struct {
	g.Meta      `path:"/finish_new_invoice" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin 完成新发票，生成支付链接或者完成自动支付"`
	InvoiceId   string `p:"invoiceId" dc:"invoiceId" v:"required#请输入InvoiceId"`
	PayMethod   int    `p:"payMethod" dc:"PayMethod,1-手动支付，2-自动支付" v:"required#请输入PayMethod"`
	DaysUtilDue int    `p:"daysUtilDue" dc:"DaysUtilDue,剩余截止支付天数" v:"required#请输入DaysUtilDue"`
}
type ProcessInvoiceForPayRes struct {
	Invoice *entity.Invoice `json:"invoice" `
}

type CancelProcessingInvoiceReq struct {
	g.Meta    `path:"/cancel_processing_invoice" tags:"Merchant-Invoice-Controller" method:"post" summary:"Admin 取消Processing 的发票"`
	InvoiceId string `p:"invoiceId" dc:"invoiceId" v:"required|length:4,30#请输入InvoiceId"`
}
type CancelProcessingInvoiceRes struct {
}
