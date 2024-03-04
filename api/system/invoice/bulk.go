package invoice

import "github.com/gogf/gf/v2/frame/g"

type BulkChannelSyncReq struct {
	g.Meta     `path:"/invoice_bulk_sync" tags:"System-Admin" method:"post" summary:"Admin Bulk Sync Invoice From Gateway (Experimental）"`
	MerchantId string `json:"merchantId" dc:"merchantId" v:"required#请输入merchantId"`
}
type BulkChannelSyncRes struct {
}

type ChannelSyncReq struct {
	g.Meta     `path:"/invoice_sync" tags:"System-Admin" method:"post" summary:"Admin Sync Invoice From Gateway (Experimental）"`
	MerchantId string `json:"merchantId" dc:"merchantId" v:"required#请输入merchantId"`
	InvoiceId  string `json:"invoiceId" dc:"invoiceId" v:"required#请输入invoiceId"`
}
type ChannelSyncRes struct {
}
