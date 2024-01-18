package balance

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/gateway/ro"
)

type DetailQueryReq struct {
	g.Meta     `path:"/merchant_balance_query" tags:"Merchant-Balance-Controller" method:"post" summary:"Merchant余额查询"`
	ChannelId  int64 `p:"channelId" dc:"channelId" v:"required#请输入 ChannelId"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
}

type DetailQueryRes struct {
	AvailableBalance       []*ro.ChannelBalance `json:"available"`
	ConnectReservedBalance []*ro.ChannelBalance `json:"connectReserved"`
	PendingBalance         []*ro.ChannelBalance `json:"pending"`
}

type UserDetailQueryReq struct {
	g.Meta     `path:"/user_balance_query" tags:"Merchant-Balance-Controller" method:"post" summary:"User余额查询"`
	UserId     int64 `p:"userId" dc:"userId" v:"required#请输入 UserId"`
	ChannelId  int64 `p:"channelId" dc:"channelId" v:"required#请输入 ChannelId"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
}

type UserDetailQueryRes struct {
	Balance              *ro.ChannelBalance   `json:"balance"`
	CashBalance          []*ro.ChannelBalance `json:"cashBalance"`
	InvoiceCreditBalance []*ro.ChannelBalance `json:"invoiceCreditBalance"`
}
