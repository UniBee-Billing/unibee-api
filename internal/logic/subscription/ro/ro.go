package ro

import entity "go-oversea-pay/internal/model/entity/oversea_pay"

type PlanDetailRo struct {
	Plan     *entity.SubscriptionPlan          `json:"plan" dc:"订阅计划"`
	Channels []*entity.SubscriptionPlanChannel `json:"channels" dc:"订阅计划 Channel 开通明细"`
	Addons   []*entity.SubscriptionPlan        `json:"addons" dc:"订阅计划 Addons 明细"`
	AddonIds []int64                           `json:"addonIds" dc:"订阅计划 Addon Ids"`
}

type SubscriptionPlanAddonParamRo struct {
	Quantity    int   `p:"quantity" dc:"数量，默认 1" `
	AddonPlanId int64 `p:"addonPlanId" dc:"订阅计划Addon ID"`
}

type SubscriptionPlanAddonRo struct {
	Quantity  int                      `p:"quantity" dc:"数量" `
	AddonPlan *entity.SubscriptionPlan `p:"addonPlan" dc:"addonPlan" `
}

type SubscriptionDetailRo struct {
	Subscription *entity.Subscription            `p:"subscription" dc:"订阅"`
	Plan         *entity.SubscriptionPlan        `p:"planId" dc:"订阅计划"`
	AddonParams  []*SubscriptionPlanAddonParamRo `p:"addonParams" dc:"订阅Addon参数"`
	Addons       []*SubscriptionPlanAddonRo      `p:"addons" dc:"订阅Addon"`
}
