package invoice

import "github.com/gogf/gf/v2/frame/g"

type BulkChannelSyncReq struct {
	g.Meta     `path:"/invoice_bulk_sync" tags:"System-Admin" method:"post" summary:"Admin Bulk Sync Invoice From Gateway (Experimental）"`
	MerchantId string `json:"merchantId" dc:"merchantId" v:"required#Require merchantId"`
}
type BulkChannelSyncRes struct {
}

type ChannelSyncReq struct {
	g.Meta     `path:"/invoice_sync" tags:"System-Admin" method:"post" summary:"Admin Sync Invoice From Gateway (Experimental）"`
	MerchantId string `json:"merchantId" dc:"merchantId" v:"required#Require merchantId"`
	InvoiceId  string `json:"invoiceId" dc:"invoiceId" v:"required#Require invoiceId"`
}
type ChannelSyncRes struct {
}
