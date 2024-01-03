package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type SubscriptionDetailReq struct {
	g.Meta         `path:"/subscription_detail" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅详情"`
	SubscriptionId int64 `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionDetailRes struct {
}

type SubscriptionChannelsReq struct {
	g.Meta     `path:"/subscription_pay_channels" tags:"User-Subscription-Controller" method:"post" summary:"订阅支持的支付渠道"`
	MerchantId int64 `p:"merchantId" d:"15621" dc:"MerchantId" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
}
type SubscriptionChannelsRes struct {
}

type SubscriptionCreateReq struct {
	g.Meta        `path:"/subscription_create" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅创建"`
	PlanId        int64                         `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	ChannelId     int64                         `p:"channelId" dc:"支付通道 ID"   v:"required#请输入 ConfirmChannelId" `
	UserId        int64                         `p:"UserId" dc:"UserId" v:"required#请输入UserId"`
	ChannelUserId string                        `p:"channelUserId" dc:"渠道用户 Id，stripe 代表 customerId" `
	Addons        []*SubscriptionPlanAddonParam `p:"addons" dc:"addons" `
}
type SubscriptionCreateRes struct {
	Subscription *entity.Subscription `json:"subscription_plan_merchant" dc:"订阅"`
}

type SubscriptionCancelReq struct {
	g.Meta         `path:"/subscription_cancel" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅取消"`
	SubscriptionId int64 `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionCancelRes struct {
}

type SubscriptionUpdateReq struct {
	g.Meta           `path:"/subscription_update" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅更新"`
	SubscriptionId   int64                         `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
	NewPlanId        int64                         `p:"newPlanId" dc:" 新的订阅计划 ID" v:"required#请输入订阅计划 ID"`
	ConfirmChannelId int64                         `p:"confirmChannelId" dc:"Web 端展示的支付通道 ID，用于验证"   v:"required#请输入 ConfirmChannelId" `
	Addons           []*SubscriptionPlanAddonParam `p:"addons" dc:"addons" `
}
type SubscriptionUpdateRes struct {
	SubscriptionPendingUpdate *entity.SubscriptionPendingUpdate `json:"subscriptionPendingUpdate" dc:"订阅"`
}

type SubscriptionPlanAddonParam struct {
	AddonPlanId int64 `p:"addonPlanId" dc:"订阅计划Addon ID" v:"required#请输入订阅计划Addon ID"`
	Quantity    int   `p:"quantity" dc:"数量，默认 1" `
}
