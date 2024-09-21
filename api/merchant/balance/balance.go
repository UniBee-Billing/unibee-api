package balance

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/gateway_bean"
)

type DetailQueryReq struct {
	g.Meta    `path:"/merchant_balance_query" tags:"Balance" method:"get" summary:"Query Merchant Gateway Balance"  deprecated:"true"`
	GatewayId uint64 `json:"gatewayId" dc:"gatewayId" v:"required"`
}

type DetailQueryRes struct {
	AvailableBalance       []*gateway_bean.GatewayBalance `json:"available"`
	ConnectReservedBalance []*gateway_bean.GatewayBalance `json:"connectReserved"`
	PendingBalance         []*gateway_bean.GatewayBalance `json:"pending"`
}

type UserDetailQueryReq struct {
	g.Meta    `path:"/user_balance_query" tags:"Balance" method:"get" summary:"Query User Balance"  deprecated:"true"`
	UserId    uint64 `json:"userId" dc:"userId" v:"required"`
	GatewayId uint64 `json:"gatewayId" dc:"gatewayId" v:"required"`
}

type UserDetailQueryRes struct {
	Balance              *gateway_bean.GatewayBalance   `json:"balance"`
	CashBalance          []*gateway_bean.GatewayBalance `json:"cashBalance"`
	InvoiceCreditBalance []*gateway_bean.GatewayBalance `json:"invoiceCreditBalance"`
}
