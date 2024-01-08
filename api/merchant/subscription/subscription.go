package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/subscription/ro"
)

type SubscriptionListReq struct {
	g.Meta     `path:"/subscription_list" tags:"Merchant-Subscription-Controller" method:"post" summary:"订阅列表"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	UserId     int64 `p:"userId"  dc:"UserId" `
	Status     int   `p:"status" dc:"不填查询所有状态，,订阅单状态，0-Init | 1-Create｜2-Active｜3-Suspend | 4-Cancel | 5-Expire" `
	Page       int   `p:"page" dc:"分页页码,0开始" `
	Count      int   `p:"count"  dc:"订阅计划货币" dc:"每页数量" `
}
type SubscriptionListRes struct {
	Subscriptions []*ro.SubscriptionDetailRo `p:"subscriptions" dc:"订阅明细"`
}

type SubscriptionCancelReq struct {
	g.Meta         `path:"/subscription_cancel_at_period_end" tags:"Merchant-Subscription-Controller" method:"post" summary:"用户订阅设置周期结束时取消"`
	SubscriptionId string `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionCancelRes struct {
}

type SubscriptionSuspendReq struct {
	g.Meta         `path:"/subscription_suspend" tags:"Merchant-Subscription-Controller" method:"post" summary:"用户订阅暂停"  deprecated:"true"`
	SubscriptionId string `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionSuspendRes struct {
}

type SubscriptionResumeReq struct {
	g.Meta         `path:"/subscription_resume" tags:"Merchant-Subscription-Controller" method:"post" summary:"用户订阅恢复"  deprecated:"true"`
	SubscriptionId string `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionResumeRes struct {
}

type SubscriptionAddNewTrialStartReq struct {
	g.Meta         `path:"/subscription_add_new_trial_start" tags:"Merchant-Subscription-Controller" method:"post" summary:"用户订阅添加试用以更改计费周期, 免费期为 currentPeriodEnd到 trailEnd 时间段"`
	SubscriptionId string `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
	TrailEnd       int64  `p:"trailEnd" dc:"新计费周期开始时间（ Unix 时间戳）-上一计费点到新周期之间为试用期，不收费" v:"required#请输入trailEnd"`
}
type SubscriptionAddNewTrialStartRes struct {
}
