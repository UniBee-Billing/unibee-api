package balance

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/gateway/ro"
)

type DetailQueryReq struct {
	g.Meta     `path:"/merchant_balance_query" tags:"Merchant-Balance-Controller" method:"post" summary:"Query Merchant Gateway Balance"`
	ChannelId  int64 `p:"channelId" dc:"channelId" v:"required"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required"`
}

type DetailQueryRes struct {
	AvailableBalance       []*ro.ChannelBalance `json:"available"`
	ConnectReservedBalance []*ro.ChannelBalance `json:"connectReserved"`
	PendingBalance         []*ro.ChannelBalance `json:"pending"`
}

type UserDetailQueryReq struct {
	g.Meta     `path:"/user_balance_query" tags:"Merchant-Balance-Controller" method:"post" summary:"Query User Balance"`
	UserId     int64 `p:"userId" dc:"userId" v:"required"`
	ChannelId  int64 `p:"channelId" dc:"channelId" v:"required"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required"`
}

type UserDetailQueryRes struct {
	Balance              *ro.ChannelBalance   `json:"balance"`
	CashBalance          []*ro.ChannelBalance `json:"cashBalance"`
	InvoiceCreditBalance []*ro.ChannelBalance `json:"invoiceCreditBalance"`
}
