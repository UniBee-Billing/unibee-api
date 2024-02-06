package refund

import "github.com/gogf/gf/v2/frame/g"

type BulkChannelSyncReq struct {
	g.Meta     `path:"/refund_bulk_sync" tags:"System-Admin-Controller" method:"post" summary:"Admin Bulk Sync Invoice Refund From Gateway (Experimental）"`
	MerchantId string `p:"merchantId" dc:"merchantId" v:"required#请输入merchantId"`
}
type BulkChannelSyncRes struct {
}
