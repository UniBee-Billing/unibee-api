package balance

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type DetailQueryReq struct {
	g.Meta    `path:"/merchant_balance_query" tags:"Merchant-Balance-Controller" method:"post" summary:"Query Merchant Gateway Balance"`
	GatewayId uint64 `p:"gatewayId" dc:"gatewayId" v:"required"`
}

type DetailQueryRes struct {
	AvailableBalance       []*ro.GatewayBalance `json:"available"`
	ConnectReservedBalance []*ro.GatewayBalance `json:"connectReserved"`
	PendingBalance         []*ro.GatewayBalance `json:"pending"`
}

type UserDetailQueryReq struct {
	g.Meta    `path:"/user_balance_query" tags:"Merchant-Balance-Controller" method:"post" summary:"Query User Balance"`
	UserId    int64  `p:"userId" dc:"userId" v:"required"`
	GatewayId uint64 `p:"gatewayId" dc:"gatewayId" v:"required"`
}

type UserDetailQueryRes struct {
	Balance              *ro.GatewayBalance   `json:"balance"`
	CashBalance          []*ro.GatewayBalance `json:"cashBalance"`
	InvoiceCreditBalance []*ro.GatewayBalance `json:"invoiceCreditBalance"`
}
