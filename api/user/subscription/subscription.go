package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/payment/gateway/ro"
	"go-oversea-pay/internal/logic/vat_gateway/base"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type SubscriptionDetailReq struct {
	g.Meta         `path:"/subscription_detail" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅详情"`
	SubscriptionId string `p:"subscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionDetailRes struct {
	Subscription               *entity.Subscription                `json:"subscription" dc:"订阅"`
	Plan                       *entity.SubscriptionPlan            `json:"planId" dc:"订阅计划"`
	Channel                    *ro.OutChannelRo                    `json:"channel" dc:"订阅渠道"`
	Addons                     []*ro.SubscriptionPlanAddonRo       `json:"addons" dc:"订阅Addon"`
	SubscriptionPendingUpdates []*entity.SubscriptionPendingUpdate `json:"subscriptionPendingUpdates" dc:"订阅更新明细"`
}

type SubscriptionPayCheckReq struct {
	g.Meta         `path:"/subscription_pay_check" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅支付状态检查"`
	SubscriptionId string `p:"subscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionPayCheckRes struct {
	PayStatus    consts.SubscriptionStatusEnum `json:"payStatus" dc:"支付状态，1-支付中，2-支付完成，3-暂停，4-取消, 5-过期"`
	Subscription *entity.Subscription          `json:"subscription" dc:"订阅"`
}

type SubscriptionChannelsReq struct {
	g.Meta     `path:"/subscription_pay_channels" tags:"User-Subscription-Controller" method:"post" summary:"订阅支持的支付渠道"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
}
type SubscriptionChannelsRes struct {
	Channels []*ro.OutChannelRo `json:"channels"`
}

type SubscriptionCreatePreviewReq struct {
	g.Meta         `path:"/subscription_create_preview" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅创建预览（仅计算）"`
	PlanId         int64                              `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	Quantity       int64                              `p:"quantity" dc:"订阅计划数量，默认 1" `
	ChannelId      int64                              `p:"channelId" dc:"支付通道 ID"   v:"required#请输入 ChannelId" `
	UserId         int64                              `p:"userId" dc:"UserId" v:"required#请输入UserId"`
	AddonParams    []*ro.SubscriptionPlanAddonParamRo `p:"addonParams" dc:"addonParams" `
	VatCountryCode string                             `p:"vatCountryCode" dc:"VatCountryCode, CountryName 缩写，Vat 接口输出"`
	VatNumber      string                             `p:"vatNumber" dc:"VatNumber, 用户输入，用于验证" `
}
type SubscriptionCreatePreviewRes struct {
	Plan              *entity.SubscriptionPlan           `json:"planId"`
	Quantity          int64                              `json:"quantity"`
	PayChannel        *ro.OutChannelRo                   `json:"payChannel"`
	AddonParams       []*ro.SubscriptionPlanAddonParamRo `json:"addonParams"`
	Addons            []*ro.SubscriptionPlanAddonRo      `json:"addons"`
	TotalAmount       int64                              `json:"totalAmount"                ` // 金额,单位：分
	Currency          string                             `json:"currency"              `      // 货币
	Invoice           *ro.ChannelDetailInvoiceRo         `json:"invoice"`
	UserId            int64                              `json:"userId" `
	Email             string                             `json:"email" `
	VatCountryCode    string                             `json:"vatCountryCode"              `
	VatCountryName    string                             `json:"vatCountryName"              `
	VatNumber         string                             `json:"vatNumber"              `
	VatNumberValidate *base.ValidResult                  `json:"vatNumberValidate"              `
}

type SubscriptionCreateReq struct {
	g.Meta             `path:"/subscription_create_submit" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅创建提交（需先调用预览接口）"`
	PlanId             int64                              `p:"planId" dc:"订阅计划 ID" v:"required#请输入订阅计划 ID"`
	Quantity           int64                              `p:"quantity" dc:"订阅计划数量，默认 1" `
	ChannelId          int64                              `p:"channelId" dc:"支付通道 ID"   v:"required#请输入 ChannelId" `
	UserId             int64                              `p:"userId" dc:"UserId" v:"required#请输入UserId"`
	AddonParams        []*ro.SubscriptionPlanAddonParamRo `p:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount int64                              `p:"confirmTotalAmount"  dc:"CreatePrepare 总金额，由Preview 接口输出"  v:"required#请输入 confirmTotalAmount"            ` // 金额,单位：分
	ConfirmCurrency    string                             `p:"confirmCurrency"  dc:"CreatePrepare 货币，由Preview 接口输出" v:"required#请输入 confirmCurrency"  `
	ReturnUrl          string                             `p:"returnUrl"  dc:"回调地址"  `
	VatCountryCode     string                             `p:"vatCountryCode" dc:"VatCountryCode, CountryName 缩写，Vat 接口输出" v:"required#请输入VatCountryCode"`
	VatNumber          string                             `p:"vatNumber" dc:"VatNumber, 用户输入，用于验证" `
}
type SubscriptionCreateRes struct {
	Subscription *entity.Subscription `json:"subscription" dc:"订阅"`
	Paid         bool                 `json:"paid"`
	Link         string               `json:"link"`
}

type SubscriptionUpdatePreviewReq struct {
	g.Meta              `path:"/subscription_update_preview" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅更新预览（仅计算）"`
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
	g.Meta              `path:"/subscription_update_submit" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅更新提交（需先调用预览接口）"`
	SubscriptionId      string                             `p:"subscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
	NewPlanId           int64                              `p:"newPlanId" dc:" 新的订阅计划 ID" v:"required#请输入订阅计划 ID"`
	Quantity            int64                              `p:"quantity" dc:"订阅计划数量，默认 1" `
	AddonParams         []*ro.SubscriptionPlanAddonParamRo `p:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount  int64                              `p:"confirmTotalAmount"  dc:"CreatePrepare 总金额，由Preview 接口输出"  v:"required#请输入 confirmTotalAmount"            ` // 金额,单位：分
	ConfirmCurrency     string                             `p:"confirmCurrency" dc:"CreatePrepare 货币，由Preview 接口输出" v:"required#请输入 confirmCurrency"  `
	ProrationDate       int64                              `p:"prorationDate" dc:"prorationDate 按比例计算开始时间，由Preview 接口输出" v:"required#请输入 prorationDate" `
	WithImmediateEffect int                                `p:"withImmediateEffect" dc:"是否立即生效，1-立即生效，2-下周期生效， withImmediateEffect=1，不会直接修改订阅，将会产生PendingUpdate 更新单和按比例发票并要求付款完成之后才会修改订阅，withImmediateEffect=2会直接修改订阅，并在下周期扣款，如果扣款失败，订阅会进入 pass_due" `
	//ConfirmChannelId int64                              `p:"confirmChannelId" dc:"Web 端展示的支付通道 ID，用于验证"   v:"required#请输入 ConfirmChannelId" `
}
type SubscriptionUpdateRes struct {
	SubscriptionPendingUpdate *entity.SubscriptionPendingUpdate `json:"subscriptionPendingUpdate" dc:"订阅"`
	Paid                      bool                              `json:"paid"`
	Link                      string                            `json:"link"`
}

type SubscriptionListReq struct {
	g.Meta     `path:"/subscription_list" tags:"User-Subscription-Controller" method:"post" summary:"订阅列表"`
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	UserId     int64  `p:"userId" dc:"UserId" v:"required|length:4,30#请输入UserId" `
	Status     int    `p:"status" dc:"不填查询所有状态，,订阅单状态，0-Init | 1-Create｜2-Active｜3-Suspend | 4-Cancel | 5-Expire" `
	SortField  string `p:"sortField" dc:"排序字段，gmt_create|gmt_modify，默认 gmt_modify" `
	SortType   string `p:"sortType" dc:"排序类型，asc|desc，默认 desc" `
	Page       int    `p:"page"  dc:"分页页码,0开始" `
	Count      int    `p:"count"  dc:"订阅计划货币" dc:"每页数量" `
}
type SubscriptionListRes struct {
	Subscriptions []*ro.SubscriptionDetailRo `p:"subscriptions" dc:"订阅明细"`
}

type SubscriptionUpdateCancelAtPeriodEndReq struct {
	g.Meta         `path:"/subscription_cancel_at_period_end" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅设置周期结束时取消"`
	SubscriptionId string `p:"subscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionUpdateCancelAtPeriodEndRes struct {
}

type SubscriptionUpdateCancelLastCancelAtPeriodEndReq struct {
	g.Meta         `path:"/subscription_cancel_last_cancel_at_period_end" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅取消上一次的周期结束时取消设置"`
	SubscriptionId string `p:"subscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionUpdateCancelLastCancelAtPeriodEndRes struct {
}

type SubscriptionSuspendReq struct {
	g.Meta         `path:"/subscription_suspend" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅暂停"  deprecated:"true"`
	SubscriptionId string `p:"subscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionSuspendRes struct {
}

type SubscriptionResumeReq struct {
	g.Meta         `path:"/subscription_resume" tags:"User-Subscription-Controller" method:"post" summary:"用户订阅恢复"  deprecated:"true"`
	SubscriptionId string `p:"subscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionResumeRes struct {
}
