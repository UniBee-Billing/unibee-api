package v1

import "github.com/gogf/gf/v2/frame/g"

type SubscriptionChannelsReq struct {
	g.Meta     `path:"/subscription_channels" tags:"Mock-Controller" method:"post" summary:"1.5订阅支付渠道"`
	MerchantId int64 `p:"merchantAccount" d:"15621" dc:"商户号" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
}
type SubscriptionChannelsRes struct {
}

type SubscriptionPlanCreateReq struct {
	g.Meta       `path:"/subscription_plan_create" tags:"Mock-Controller" method:"post" summary:"1.6订阅计划创建"`
	MerchantId   int64  `p:"merchantAccount" d:"15621" dc:"商户号" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
	PlanName     string `p:"planName"    v:"required|length:4,30#请输入订阅计划名称长度为:{min}到:{max}位" ` // 计划名称
	Amount       int64  `p:"amount"      v:"required#请输入订阅计划金额" `                              // 金额,单位：分
	Currency     string `p:"currency"    v:"required#请输入订阅计划货币" `                              // 货币
	IntervalUnit string `p:"intervalUnit" v:"required#请输入订阅计划周期，小写: day|month|year|week" `     // 周期,day|month|year|week
	Description  string `p:"description"  `                                                    //
	ImageUrl     string `p:"imageUrl"     `                                                    // image_url
	HomeUrl      string `p:"homeUrl"      `                                                    // home_url
}
type SubscriptionPlanCreateRes struct {
	//PlanId int64 `json:"planId"      dc:"订阅计划 ID" `
}

type SubscriptionPlanChannelTransferReq struct {
	g.Meta    `path:"/subscription_plan_channel_transfer" tags:"Mock-Controller" method:"post" summary:"1.7订阅计划支付通道开通"`
	PlanId    int64  `p:"planId" d:"15621" dc:"订阅计划 ID" v:"required|请输入订阅计划 ID"`
	ChannelId string `p:"channelId"    v:"required|请输入 ChannelId" `
}
type SubscriptionPlanChannelTransferRes struct {
}

type SubscriptionPlanDetailReq struct {
	g.Meta `path:"/subscription_plan_detail" tags:"Mock-Controller" method:"post" summary:"1.8订阅计划明细"`
	PlanId int64 `p:"planId" d:"15621" dc:"订阅计划 ID" v:"required|请输入订阅计划 ID"`
}
type SubscriptionPlanDetailRes struct {
}
