package payment

import "github.com/gogf/gf/v2/frame/g"

type BulkChannelSyncReq struct {
	g.Meta     `path:"/invoice_payment_bulk_sync" tags:"System-Admin-Controller" method:"post" summary:"Admin Bulk Sync Invoice Payment From Gateway"`
	MerchantId string `p:"merchantId" dc:"merchantId" v:"required#请输入merchantId"`
}
type BulkChannelSyncRes struct {
}
