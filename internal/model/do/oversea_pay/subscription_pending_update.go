// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionPendingUpdate is the golang structure of table subscription_pending_update for DAO operations like Where/Data.
type SubscriptionPendingUpdate struct {
	g.Meta               `orm:"table:subscription_pending_update, do:true"`
	Id                   interface{} //
	MerchantId           interface{} // 商户Id
	SubscriptionId       interface{} // 订阅id（内部编号）
	UpdateSubscriptionId interface{} // 升级单ID（内部编号）
	ChannelUpdateId      interface{} // 支付渠道订阅更新单id， stripe 适用 channelInvoiceId对应
	GmtCreate            *gtime.Time // 创建时间
	Amount               interface{} // 金额,单位：分
	Status               interface{} // 订阅单状态，0-Init | 1-Create｜2-Finished｜3-Cancelled
	UpdateAmount         interface{} // 升级到金额,单位：分
	Currency             interface{} // 货币
	UpdateCurrency       interface{} // 升级到货币
	PlanId               interface{} // 计划ID
	UpdatePlanId         interface{} // 升级到计划ID
	Quantity             interface{} // quantity
	UpdateQuantity       interface{} // 升级到quantity
	AddonData            interface{} // plan addon json data
	UpdateAddonData      interface{} // 升级到plan addon json data
	ChannelId            interface{} // 支付渠道Id
	UserId               interface{} // userId
	GmtModify            *gtime.Time // 修改时间
	IsDeleted            interface{} //
	Paid                 interface{} // 是否已支付，0-否，1-是
	Link                 interface{} // 支付链接
	ChannelStatus        interface{} // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	MerchantUserId       interface{} // merchant_user_id
	Data                 interface{} // 渠道额外参数，JSON格式
	ResponseData         interface{} // 渠道返回参数，JSON格式
	EffectImmediate      interface{} // 是否马上生效，0-否，1-是
	EffectTime           interface{} // effect_immediate=0, 预计生效时间 unit_time
	AdminNote            interface{} // Admin 修改备注
}
