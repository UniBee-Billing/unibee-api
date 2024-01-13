package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type SubscriptionDetailReq struct {
	g.Meta         `path:"/subscription_detail" tags:"Merchant-Subscription-Controller" method:"post" summary:"订阅详情"`
	SubscriptionId string `p:"subscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionDetailRes struct {
	Subscription               *entity.Subscription                `json:"subscription" dc:"订阅"`
	Plan                       *entity.SubscriptionPlan            `json:"planId" dc:"订阅计划"`
	Channel                    *ro.OutChannelRo                    `json:"channel" dc:"订阅渠道"`
	Addons                     []*ro.SubscriptionPlanAddonRo       `json:"addons" dc:"订阅Addon"`
	SubscriptionPendingUpdates []*entity.SubscriptionPendingUpdate `json:"subscriptionPendingUpdates" dc:"订阅更新明细"`
}

type SubscriptionListReq struct {
	g.Meta     `path:"/subscription_list" tags:"Merchant-Subscription-Controller" method:"post" summary:"订阅列表"`
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	UserId     int64  `p:"userId"  dc:"UserId" `
	Status     int    `p:"status" dc:"不填查询所有状态，,订阅单状态，0-Init | 1-Create｜2-Active｜3-Suspend | 4-Cancel | 5-Expire" `
	SortField  string `p:"sortField" dc:"排序字段，gmt_create|gmt_modify，默认 gmt_modify" `
	SortType   string `p:"sortType" dc:"排序类型，asc|desc，默认 desc" `
	Page       int    `p:"page" dc:"分页页码,0开始" `
	Count      int    `p:"count"  dc:"订阅计划货币" dc:"每页数量" `
}
type SubscriptionListRes struct {
	Subscriptions []*ro.SubscriptionDetailRo `json:"subscriptions" dc:"订阅明细"`
}

type SubscriptionCancelReq struct {
	g.Meta         `path:"/subscription_cancel_at_period_end" tags:"Merchant-Subscription-Controller" method:"post" summary:"Merchant 修改用户订阅-设置周期结束时取消"`
	SubscriptionId string `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionCancelRes struct {
}

type SubscriptionSuspendReq struct {
	g.Meta         `path:"/subscription_suspend" tags:"Merchant-Subscription-Controller" method:"post" summary:"Merchant 修改用户订阅-暂停"  deprecated:"true"`
	SubscriptionId string `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionSuspendRes struct {
}

type SubscriptionResumeReq struct {
	g.Meta         `path:"/subscription_resume" tags:"Merchant-Subscription-Controller" method:"post" summary:"Merchant 修改用户订阅-恢复"  deprecated:"true"`
	SubscriptionId string `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionResumeRes struct {
}

type SubscriptionAddNewTrialStartReq struct {
	g.Meta         `path:"/subscription_add_new_trial_start" tags:"Merchant-Subscription-Controller" method:"post" summary:"Merchant 修改用户订阅-添加试用以更改计费周期, 免费期为 currentPeriodEnd到 trailEnd 时间段"`
	SubscriptionId string `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
	TrailEnd       int64  `p:"trailEnd" dc:"新计费周期开始时间（ Unix 时间戳）-上一计费点到新周期之间为试用期，不收费" v:"required#请输入trailEnd"`
}
type SubscriptionAddNewTrialStartRes struct {
}

type SubscriptionUpdatePreviewReq struct {
	g.Meta              `path:"/subscription_update_preview" tags:"Merchant-Subscription-Controller" method:"post" summary:"Merchant 修改用户订阅-更新预览（仅计算）"`
	SubscriptionId      string                             `p:"subscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
	NewPlanId           int64                              `p:"newPlanId" dc:" 新的订阅计划 ID" v:"required#请输入订阅计划 ID"`
	Quantity            int64                              `p:"quantity" dc:"订阅计划数量，默认 1" `
	WithImmediateEffect int                                `p:"withImmediateEffect" dc:"是否立即生效，1-立即生效，2-下周期生效, withImmediateEffect=1，不会直接修改订阅，将会产生PendingUpdate 更新单和按比例发票并要求付款完成之后才会修改订阅，withImmediateEffect=2会直接修改订阅，并在下周期扣款，如果扣款失败，订阅会进入 pass_due" `
	AddonParams         []*ro.SubscriptionPlanAddonParamRo `p:"addonParams" dc:"addonParams" `
}
type SubscriptionUpdatePreviewRes struct {
	TotalAmount   int64                      `json:"totalAmount"                ` // 金额,单位：分
	Currency      string                     `json:"currency"              `      // 货币
	Invoice       *ro.ChannelDetailInvoiceRo `json:"invoice"`
	ProrationDate int64                      `json:"prorationDate"`
}

type SubscriptionUpdateReq struct {
	g.Meta              `path:"/subscription_update_submit" tags:"Merchant-Subscription-Controller" method:"post" summary:"Merchant 修改用户订阅-更新提交（需先调用预览接口）"`
	SubscriptionId      string                             `p:"subscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
	NewPlanId           int64                              `p:"newPlanId" dc:" 新的订阅计划 ID" v:"required#请输入订阅计划 ID"`
	Quantity            int64                              `p:"quantity" dc:"订阅计划数量，默认 1" `
	AddonParams         []*ro.SubscriptionPlanAddonParamRo `p:"addonParams" dc:"addonParams" `
	WithImmediateEffect int                                `p:"withImmediateEffect" dc:"是否立即生效，1-立即生效，2-下周期生效, withImmediateEffect=1，不会直接修改订阅，将会产生PendingUpdate 更新单和按比例发票并要求付款完成之后才会修改订阅，withImmediateEffect=2会直接修改订阅，并在下周期扣款，如果扣款失败，订阅会进入 pass_due" `
	ConfirmTotalAmount  int64                              `p:"confirmTotalAmount"  dc:"CreatePrepare 总金额，由Preview 接口输出"  v:"required#请输入 confirmTotalAmount"            ` // 金额,单位：分
	ConfirmCurrency     string                             `p:"confirmCurrency" dc:"CreatePrepare 货币，由Preview 接口输出" v:"required#请输入 confirmCurrency"  `
	ProrationDate       int64                              `p:"prorationDate" dc:"prorationDate 按比例计算开始时间，由Preview 接口输出" v:"required#请输入 prorationDate" `
	//ConfirmChannelId int64                              `p:"confirmChannelId" dc:"Web 端展示的支付通道 ID，用于验证"   v:"required#请输入 ConfirmChannelId" `
}

type SubscriptionUpdateRes struct {
	SubscriptionPendingUpdate *entity.SubscriptionPendingUpdate `json:"subscriptionPendingUpdate" dc:"订阅"`
	Paid                      bool                              `json:"paid"`
	Link                      string                            `json:"link"`
}
