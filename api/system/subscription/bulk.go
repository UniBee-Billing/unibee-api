package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
)

type BulkChannelSyncReq struct {
	g.Meta     `path:"/sub_bulk_sync" tags:"System-Admin-Controller" method:"post" summary:"Admin Bulk Sync Subscription From Gateway"`
	MerchantId string `p:"merchantId" dc:"merchantId" v:"required#请输入merchantId"`
}
type BulkChannelSyncRes struct {
}
