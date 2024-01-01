package ro

import entity "go-oversea-pay/internal/model/entity/oversea_pay"

type SubscriptionPlanRo struct {
	Plan     *entity.SubscriptionPlan          `json:"plan" dc:"订阅计划"`
	Channels []*entity.SubscriptionPlanChannel `json:"channels" dc:"订阅计划 Channel 开通明细"`
	Addons   []*entity.SubscriptionPlan        `json:"addons" dc:"订阅计划 Addons 明细"`
	AddonIds []int64                           `json:"addonIds" dc:"订阅计划 Addon Ids"`
}
