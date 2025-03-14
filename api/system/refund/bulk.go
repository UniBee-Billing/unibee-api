package refund

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type BulkChannelSyncReq struct {
	g.Meta     `path:"/refund_bulk_sync" tags:"System-Admin" method:"post" summary:"Admin Bulk Sync Invoice Refund From Gateway (Experimental）"`
	MerchantId string `json:"merchantId" dc:"merchantId" v:"required#Require merchantId"`
}
type BulkChannelSyncRes struct {
}

type GatewayDetailReq struct {
	g.Meta   `path:"/refund_gateway_detail" tags:"System-Admin" method:"post" summary:"Admin Trigger Refund Gateway Detail"`
	RefundId string `json:"refundId" v:"required#Require refundId"`
}

type GatewayDetailRes struct {
	RefundDetail *gjson.Json `json:"refundDetail"`
}
