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
	GmtCreate            *gtime.Time // 创建时间
	Amount               interface{} // 金额,单位：分
	UpdateAmount         interface{} // 升级到金额,单位：分
	Currency             interface{} // 货币
	UpdateCurrency       interface{} // 升级到货币
	PlanId               interface{} // 计划ID
	UpdatePlanId         interface{} // 升级到计划ID
	Quantity             interface{} // quantity
	UpdateQuantity       interface{} // 升级到quantity
	AddonData            interface{} // plan addon json data
	UpdatedAddonData     interface{} // 升级到plan addon json data
	ChannelId            interface{} // 支付渠道Id
	Status               interface{} // 订阅单状态，0-Init | 1-Create｜2-Active｜3-Suspend
	UserId               interface{} // userId
	ChannelUpdateId      interface{} // 支付渠道订阅更新单id
	Data                 interface{} // 渠道额外参数，JSON格式
	ResponseData         interface{} // 渠道返回参数，JSON格式
	GmtModify            *gtime.Time // 修改时间
	IsDeleted            interface{} //
	Link                 interface{} //
	ChannelStatus        interface{} // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
}
