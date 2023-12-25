package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type SubscriptionChannelsReq struct {
	g.Meta     `path:"/subscription_pay_channels" tags:"Subscription-Controller" method:"post" summary:"1.1订阅支持的支付渠道"`
	MerchantId int64 `p:"merchantAccount" d:"15621" dc:"商户号" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
}
type SubscriptionChannelsRes struct {
}

type SubscriptionPlanCreateReq struct {
	g.Meta             `path:"/subscription_plan_create_and_activate" tags:"Subscription-Plan-Controller" method:"post" summary:"1.2订阅计划创建"`
	MerchantId         int64  `p:"merchantAccount" d:"15621" dc:"商户号" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
	PlanName           string `p:"planName"    v:"required|length:4,30#请输入订阅计划名称长度为:{min}到:{max}位" ` // 计划名称
	Amount             int64  `p:"amount"      v:"required#请输入订阅计划金额" `                              // 金额,单位：分
	Currency           string `p:"currency"    v:"required#请输入订阅计划货币" `                              // 货币
	IntervalUnit       string `p:"intervalUnit" v:"required#请输入订阅计划周期，小写: day|month|year|week" `     // 周期,day|month|year|week
	Description        string `p:"description"  `                                                    //
	ProductName        string `p:"productName" dc:"不填默认 PlanName"  `                                 //
	ProductDescription string `p:"productDescription" dc:"不填默认 Description" `                        //
	ImageUrl           string `p:"imageUrl"     v:"required#请输入ImageUrl,需 http 开头" `                 // image_url
	HomeUrl            string `p:"homeUrl"      `                                                    // home_url
}
type SubscriptionPlanCreateRes struct {
	Plan *entity.SubscriptionPlan `json:"plan" dc:"订阅计划"`
}

type SubscriptionPlanChannelTransferAndActivateReq struct {
	g.Meta    `path:"/subscription_plan_channel_transfer_and_activate" tags:"Subscription-Plan-Controller" method:"post" summary:"1.3订阅计划支付通道开通并激活"`
	PlanId    int64 `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	ChannelId int64 `p:"channelId"    v:"required#请输入 ChannelId" `
}
type SubscriptionPlanChannelTransferAndActivateRes struct {
}

type SubscriptionPlanChannelActivateReq struct {
	g.Meta    `path:"/subscription_plan_channel_activate" tags:"Subscription-Plan-Controller" method:"post" summary:"1.4订阅计划支付通道激活"`
	PlanId    int64 `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	ChannelId int64 `p:"channelId"    v:"required#请输入 ChannelId" `
}
type SubscriptionPlanChannelActivateRes struct {
}

type SubscriptionPlanChannelDeactivateReq struct {
	g.Meta    `path:"/subscription_plan_channel_deactivate" tags:"Subscription-Plan-Controller" method:"post" summary:"1.5订阅计划支付通道取消激活"`
	PlanId    int64 `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	ChannelId int64 `p:"channelId"    v:"required#请输入 ChannelId" `
}
type SubscriptionPlanChannelDeactivateRes struct {
}

type SubscriptionPlanDetailReq struct {
	g.Meta `path:"/subscription_plan_detail" tags:"Subscription-Plan-Controller" method:"post" summary:"1.6订阅计划详情"`
	PlanId int64 `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
}
type SubscriptionPlanDetailRes struct {
	Plan     *entity.SubscriptionPlan          `json:"plan" dc:"订阅计划"`
	Channels *[]entity.SubscriptionPlanChannel `json:"channels" dc:"订阅计划 Channel 开通明细"`
}

type SubscriptionCreateReq struct {
	g.Meta        `path:"/subscription_create" tags:"Subscription-Controller" method:"post" summary:"1.7用户订阅创建"`
	PlanId        int64  `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	ChannelId     int64  `p:"channelId" dc:"支付通道 ID"   v:"required#请输入 ChannelId" `
	UserId        int64  `p:"UserId" dc:"UserId" v:"required#请输入UserId"`
	ChannelUserId string `p:"channelUserId" dc:"渠道用户 Id，stripe 代表 customerId" `
}
type SubscriptionCreateRes struct {
	Subscription *entity.Subscription `json:"subscription" dc:"订阅"`
}

type SubscriptionCancelReq struct {
	g.Meta         `path:"/subscription_cancel" tags:"Subscription-Controller" method:"post" summary:"1.8用户订阅取消"`
	SubscriptionId int64 `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionCancelRes struct {
}

type SubscriptionDetailReq struct {
	g.Meta         `path:"/subscription_detail" tags:"Subscription-Controller" method:"post" summary:"1.9用户订阅详情"`
	SubscriptionId int64 `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionDetailRes struct {
}

type SubscriptionUpdateReq struct {
	g.Meta         `path:"/subscription_update" tags:"Subscription-Controller" method:"post" summary:"1.10用户订阅更新"`
	SubscriptionId int64 `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
	NewPlanId      int64 `p:"newPlanId" dc:" 新的订阅计划 ID" v:"required#请输入订阅计划 ID"`
	ChannelId      int64 `p:"channelId" dc:"支付通道 ID"   v:"required#请输入 ChannelId" `
}
type SubscriptionUpdateRes struct {
	Subscription *entity.Subscription `json:"subscription" dc:"订阅"`
}
