package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/subscription/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type SubscriptionChannelsReq struct {
	g.Meta     `path:"/subscription_pay_channels" tags:"Subscription-Controller" method:"post" summary:"订阅支持的支付渠道"`
	MerchantId int64 `p:"merchantId" d:"15621" dc:"MerchantId" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
}
type SubscriptionChannelsRes struct {
}

type SubscriptionPlanCreateReq struct {
	g.Meta             `path:"/subscription_plan_create_and_activate" tags:"Subscription-Plan-Controller" method:"post" summary:"订阅计划创建"`
	MerchantId         int64  `p:"merchantId" d:"15621" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	PlanName           string `p:"planName" dc:"订阅计划名称"   v:"required|length:4,30#请输入订阅计划名称长度为:{min}到:{max}位" `                                                       // 计划名称
	Amount             int64  `p:"amount"   dc:"订阅计划金额"   v:"required#请输入订阅计划金额" `                                                                                    // 金额,单位：分
	Currency           string `p:"currency"   dc:"订阅计划货币" v:"required#请输入订阅计划货币" `                                                                                    // 货币
	IntervalUnit       string `p:"intervalUnit" dc:"订阅计划周期，小写: day|month|year|week" v:"required#请输入订阅计划周期，小写: day|month|year|week" `                                  // 周期,day|month|year|week
	IntervalCount      int    `p:"intervalCount"  d:"1" dc:"不输入或者输入值小于 1，强制设置 1，订阅计费之间的间隔数。例如，每 3 个月interval=month计费一次interval_count=3。允许的最长间隔为一年（1 年、12 个月或 52 周）" ` // 金额,单位：分
	Type               int    `p:"type"  d:"1"  dc:"默认值 1，,1-main plan，2-addon plan" `                                                                                // 金额,单位：分
	Description        string `p:"description"  dc:"描述"`                                                                                                              //
	ProductName        string `p:"productName" dc:"不填默认 PlanName"  `                                                                                                  //
	ProductDescription string `p:"productDescription" dc:"不填默认 Description" `                                                                                         //
	ImageUrl           string `p:"imageUrl"    dc:"ImageUrl,需 http 开头"  v:"required#请输入ImageUrl,需 http 开头" `                                                          // image_url
	HomeUrl            string `p:"homeUrl"    dc:"HomeUrl,需 http 开头"  `                                                                                               // home_url
}
type SubscriptionPlanCreateRes struct {
	Plan *entity.SubscriptionPlan `json:"plan" dc:"订阅计划"`
}

type SubscriptionPlanEditReq struct {
	g.Meta             `path:"/subscription_plan_edit" tags:"Subscription-Plan-Controller" method:"post" summary:"订阅计划修改(覆盖模式）"`
	PlanId             int64  `p:"planId" dc:"PlanId" v:"required#请输入订阅计划 ID"`
	PlanName           string `p:"planName" dc:"订阅计划名称"   v:"required|length:4,30#请输入订阅计划名称长度为:{min}到:{max}位" `                                                       // 计划名称
	Amount             int64  `p:"amount"   dc:"订阅计划金额"   v:"required#请输入订阅计划金额" `                                                                                    // 金额,单位：分
	Currency           string `p:"currency"   dc:"订阅计划货币" v:"required#请输入订阅计划货币" `                                                                                    // 货币
	IntervalUnit       string `p:"intervalUnit" dc:"订阅计划周期，小写: day|month|year|week" v:"required#请输入订阅计划周期，小写: day|month|year|week" `                                  // 周期,day|month|year|week
	IntervalCount      int    `p:"intervalCount"  d:"1" dc:"不输入或者输入值小于 1，强制设置 1，订阅计费之间的间隔数。例如，每 3 个月interval=month计费一次interval_count=3。允许的最长间隔为一年（1 年、12 个月或 52 周）" ` // 金额,单位：分
	Description        string `p:"description"  dc:"描述"`                                                                                                              //
	ProductName        string `p:"productName" dc:"不填默认 PlanName"  `                                                                                                  //
	ProductDescription string `p:"productDescription" dc:"不填默认 Description" `                                                                                         //
	ImageUrl           string `p:"imageUrl"    dc:"ImageUrl,需 http 开头"  v:"required#请输入ImageUrl,需 http 开头" `                                                          // image_url
	HomeUrl            string `p:"homeUrl"    dc:"HomeUrl,需 http 开头"  `                                                                                               // home_url
}
type SubscriptionPlanEditRes struct {
	Plan *entity.SubscriptionPlan `json:"plan" dc:"订阅计划"`
}

type SubscriptionPlanAddonsBindingReq struct {
	g.Meta   `path:"/subscription_plan_addons_binding" tags:"Subscription-Plan-Controller" method:"post" summary:"订阅计划 Addons 绑定"`
	PlanId   int64   `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	Action   int64   `p:"action" d:"0" dc:"操作类型，0-覆盖,1-添加，2-删除" v:"required#请输入操作类型"`
	AddonIds []int64 `p:"addonIds"  dc:"addon 类型 Plan Ids"  v:"required#请输入 addonIds" `
}
type SubscriptionPlanAddonsBindingRes struct {
	Plan *entity.SubscriptionPlan `json:"plan" dc:"订阅计划"`
}

type SubscriptionPlanListReq struct {
	g.Meta     `path:"/subscription_plan_list" tags:"Subscription-Plan-Controller" method:"post" summary:"订阅计划列表"`
	MerchantId int64  `p:"merchantId" d:"15621" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	Type       int    `p:"type"  d:"1"  dc:"不填查询所有类型，,1-main plan，2-addon plan" `
	Status     int    `p:"status" dc:"不填查询所有状态，,状态，1-编辑中，2-活跃，3-非活跃，4-过期" `
	Currency   string `p:"currency" d:"usd"  dc:"订阅计划货币"  `
	Page       int    `p:"page" d:"0"  dc:"分页页码,0开始" `
	Count      int    `p:"count" d:"20"  dc:"订阅计划货币" dc:"每页数量" `
}
type SubscriptionPlanListRes struct {
	Plans []*ro.SubscriptionPlanRo `p:"plans" dc:"订阅计划明细"`
}

type SubscriptionPlanChannelTransferAndActivateReq struct {
	g.Meta `path:"/subscription_plan_activate" tags:"Subscription-Plan-Controller" method:"post" summary:"订阅计划支付通道激活并发布"`
	PlanId int64 `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	//ChannelId int64 `p:"channelId"    v:"required#请输入 ConfirmChannelId" `
}
type SubscriptionPlanChannelTransferAndActivateRes struct {
}

type SubscriptionPlanChannelActivateReq struct {
	g.Meta    `path:"/subscription_plan_channel_activate" tags:"Subscription-Plan-Controller" method:"post" summary:"订阅计划支付单通道激活"`
	PlanId    int64 `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	ChannelId int64 `p:"channelId"    v:"required#请输入 ConfirmChannelId" `
}
type SubscriptionPlanChannelActivateRes struct {
}

type SubscriptionPlanChannelDeactivateReq struct {
	g.Meta    `path:"/subscription_plan_channel_deactivate" tags:"Subscription-Plan-Controller" method:"post" summary:"订阅计划支付单通道取消激活"`
	PlanId    int64 `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	ChannelId int64 `p:"channelId"    v:"required#请输入 ConfirmChannelId" `
}
type SubscriptionPlanChannelDeactivateRes struct {
}

type SubscriptionPlanDetailReq struct {
	g.Meta `path:"/subscription_plan_detail" tags:"Subscription-Plan-Controller" method:"post" summary:"订阅计划详情"`
	PlanId int64 `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
}
type SubscriptionPlanDetailRes struct {
	Plan *ro.SubscriptionPlanRo `p:"plan" dc:"订阅计划明细"`
}

type SubscriptionPlanExpireReq struct {
	g.Meta    `path:"/subscription_plan_expire" tags:"Subscription-Plan-Controller" method:"post" summary:"订阅计划过期"`
	PlanId    int64 `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	EmailCode int64 `p:"emailCode" dc:"邮箱中获取的验证码" v:"required#请输入验证码"`
}
type SubscriptionPlanExpireRes struct {
}

type SubscriptionCreateReq struct {
	g.Meta        `path:"/subscription_create" tags:"Subscription-Controller" method:"post" summary:"用户订阅创建"`
	PlanId        int64  `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	ChannelId     int64  `p:"channelId" dc:"支付通道 ID"   v:"required#请输入 ConfirmChannelId" `
	UserId        int64  `p:"UserId" dc:"UserId" v:"required#请输入UserId"`
	ChannelUserId string `p:"channelUserId" dc:"渠道用户 Id，stripe 代表 customerId" `
}
type SubscriptionCreateRes struct {
	Subscription *entity.Subscription `json:"subscription" dc:"订阅"`
}

type SubscriptionCancelReq struct {
	g.Meta         `path:"/subscription_cancel" tags:"Subscription-Controller" method:"post" summary:"用户订阅取消"`
	SubscriptionId int64 `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionCancelRes struct {
}

type SubscriptionDetailReq struct {
	g.Meta         `path:"/subscription_detail" tags:"Subscription-Controller" method:"post" summary:"用户订阅详情"`
	SubscriptionId int64 `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionDetailRes struct {
}

type SubscriptionUpdateReq struct {
	g.Meta           `path:"/subscription_update" tags:"Subscription-Controller" method:"post" summary:"用户订阅更新"`
	SubscriptionId   int64 `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
	NewPlanId        int64 `p:"newPlanId" dc:" 新的订阅计划 ID" v:"required#请输入订阅计划 ID"`
	ConfirmChannelId int64 `p:"confirmChannelId" dc:"Web 端展示的支付通道 ID，用于验证"   v:"required#请输入 ConfirmChannelId" `
}
type SubscriptionUpdateRes struct {
	Subscription *entity.Subscription `json:"subscription" dc:"订阅"`
}

type SubscriptionWebhookCheckAndSetupReq struct {
	g.Meta `path:"/subscription_webhook_check_and_setup" tags:"Subscription-Controller" method:"post" summary:"Webhook 初始化"`
}
type SubscriptionWebhookCheckAndSetupRes struct {
}
