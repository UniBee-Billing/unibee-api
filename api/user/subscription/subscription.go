package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/subscription/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type SubscriptionDetailReq struct {
	g.Meta         `path:"/subscription_detail" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅详情"`
	SubscriptionId string `p:"ubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionDetailRes struct {
	Subscription              *entity.Subscription                `p:"subscription" dc:"订阅"`
	Plan                      *entity.SubscriptionPlan            `p:"planId" dc:"订阅计划"`
	Addons                    []*ro.SubscriptionPlanAddonRo       `p:"addons" dc:"订阅Addon"`
	SubscriptionPendingUpdate []*entity.SubscriptionPendingUpdate `p:"subscriptionPendingUpdate" dc:"订阅更新明细"`
}

type SubscriptionChannelsReq struct {
	g.Meta     `path:"/subscription_pay_channels" tags:"User-Subscription-Controller" method:"post" summary:"订阅支持的支付渠道"`
	MerchantId int64 `p:"merchantId" d:"15621" dc:"MerchantId" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
}
type SubscriptionChannelsRes struct {
}

type SubscriptionCreatePreviewReq struct {
	g.Meta      `path:"/subscription_create_preview" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅创建计算（仅计算）"`
	PlanId      int64                              `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	Quantity    int64                              `p:"quantity" dc:"订阅计划数量，默认 1" `
	ChannelId   int64                              `p:"channelId" dc:"支付通道 ID"   v:"required#请输入 ConfirmChannelId" `
	UserId      int64                              `p:"UserId" dc:"UserId" v:"required#请输入UserId"`
	AddonParams []*ro.SubscriptionPlanAddonParamRo `p:"addonParams" dc:"addonParams" `
}
type SubscriptionCreatePreviewRes struct {
	Plan        *entity.SubscriptionPlan           `json:"planId"`
	Quantity    int64                              `json:"quantity"`
	PayChannel  *entity.OverseaPayChannel          `json:"payChannel"`
	AddonParams []*ro.SubscriptionPlanAddonParamRo `json:"addonParams"`
	Addons      []*ro.SubscriptionPlanAddonRo      `json:"addons"`
	TotalAmount int64                              `json:"totalAmount"                ` // 金额,单位：分
	Currency    string                             `json:"currency"              `      // 货币
	Invoice     *ro.SubscriptionInvoiceRo          `json:"invoice"`
	UserId      int64                              `json:"userId" `
	Email       string                             `json:"email" `
}

type SubscriptionCreateReq struct {
	g.Meta             `path:"/subscription_create_submit" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅创建提交"`
	PlanId             int64                              `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	Quantity           int64                              `p:"quantity" dc:"订阅计划数量，默认 1" `
	ChannelId          int64                              `p:"channelId" dc:"支付通道 ID"   v:"required#请输入 ConfirmChannelId" `
	UserId             int64                              `p:"UserId" dc:"UserId" v:"required#请输入UserId"`
	AddonParams        []*ro.SubscriptionPlanAddonParamRo `p:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount int64                              `p:"confirmTotalAmount"  dc:"CreatePrepare 接口输出的总金额"  v:"required#请输入 confirmTotalAmount"            ` // 金额,单位：分
	ConfirmCurrency    string                             `p:"confirmCurrency" d:"usd"  dc:"CreatePrepare 接口输出的货币" v:"required#请输入 confirmCurrency"  `
}
type SubscriptionCreateRes struct {
	Subscription *entity.Subscription      `json:"subscription" dc:"订阅"`
	Invoice      *ro.SubscriptionInvoiceRo `json:"invoice"`
}

type SubscriptionCancelReq struct {
	g.Meta         `path:"/subscription_cancel" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅取消"`
	SubscriptionId string `p:"SubscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionCancelRes struct {
}

type SubscriptionUpdatePreviewReq struct {
	g.Meta         `path:"/subscription_update_preview" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅更新计算（仅计算）"`
	SubscriptionId string                             `p:"subscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
	NewPlanId      int64                              `p:"newPlanId" dc:" 新的订阅计划 ID" v:"required#请输入订阅计划 ID"`
	Quantity       int64                              `p:"quantity" dc:"订阅计划数量，默认 1" `
	AddonParams    []*ro.SubscriptionPlanAddonParamRo `p:"addonParams" dc:"addonParams" `
}
type SubscriptionUpdatePreviewRes struct {
	TotalAmount   int64                     `json:"totalAmount"                ` // 金额,单位：分
	Currency      string                    `json:"currency"              `      // 货币
	Invoice       *ro.SubscriptionInvoiceRo `json:"invoice"`
	ProrationDate int64                     `json:"prorationDate"`
}

type SubscriptionUpdateReq struct {
	g.Meta             `path:"/subscription_update_submit" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅更新提交"`
	SubscriptionId     string                             `p:"subscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
	NewPlanId          int64                              `p:"newPlanId" dc:" 新的订阅计划 ID" v:"required#请输入订阅计划 ID"`
	Quantity           int64                              `p:"quantity" dc:"订阅计划数量，默认 1" `
	AddonParams        []*ro.SubscriptionPlanAddonParamRo `p:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount int64                              `p:"confirmTotalAmount"  dc:"CreatePrepare 接口输出的总金额"  v:"required#请输入 confirmTotalAmount"            ` // 金额,单位：分
	ConfirmCurrency    string                             `p:"confirmCurrency" d:"usd"  dc:"CreatePrepare 接口输出的货币" v:"required#请输入 confirmCurrency"  `
	ProrationDate      int64                              `p:"prorationDate" dc:"prorationDate 按比例计算开始时间，由Preview 接口输出" v:"required#请输入 prorationDate" `
	//ConfirmChannelId int64                              `p:"confirmChannelId" dc:"Web 端展示的支付通道 ID，用于验证"   v:"required#请输入 ConfirmChannelId" `
}
type SubscriptionUpdateRes struct {
	SubscriptionPendingUpdate *entity.SubscriptionPendingUpdate `json:"subscriptionPendingUpdate" dc:"订阅"`
	TotalAmount               int64                             `json:"totalAmount"                ` // 金额,单位：分
	Currency                  string                            `json:"currency"              `      // 货币
	Invoice                   *ro.SubscriptionInvoiceRo         `json:"invoice"`
}

type SubscriptionListReq struct {
	g.Meta     `path:"/subscription_list" tags:"User-Subscription-Controller" method:"post" summary:"订阅列表"`
	MerchantId int64 `p:"merchantId" d:"15621" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	UserId     int64 `p:"userId"  d:"1" dc:"UserId" v:"required|length:4,30#请输入UserId" `
	Status     int   `p:"status" dc:"不填查询所有状态，,订阅单状态，0-Init | 1-Create｜2-Active｜3-Suspend | 4-Cancel | 5-Expire" `
	Page       int   `p:"page" d:"0"  dc:"分页页码,0开始" `
	Count      int   `p:"count" d:"20"  dc:"订阅计划货币" dc:"每页数量" `
}
type SubscriptionListRes struct {
	Subscriptions []*ro.SubscriptionDetailRo `p:"subscriptions" dc:"订阅明细"`
}
