// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionPendingUpdate is the golang structure for table subscription_pending_update.
type SubscriptionPendingUpdate struct {
	Id                   uint64      `json:"id"                   description:""`                                                                                                                                                       //
	MerchantId           int64       `json:"merchantId"           description:"商户Id"`                                                                                                                                                   // 商户Id
	SubscriptionId       string      `json:"subscriptionId"       description:"订阅id（内部编号）"`                                                                                                                                             // 订阅id（内部编号）
	UpdateSubscriptionId string      `json:"updateSubscriptionId" description:"升级单ID（内部编号）"`                                                                                                                                            // 升级单ID（内部编号）
	ChannelUpdateId      string      `json:"channelUpdateId"      description:"支付渠道订阅更新单id， stripe 适用 channelInvoiceId对应"`                                                                                                              // 支付渠道订阅更新单id， stripe 适用 channelInvoiceId对应
	GmtCreate            *gtime.Time `json:"gmtCreate"            description:"创建时间"`                                                                                                                                                   // 创建时间
	Amount               int64       `json:"amount"               description:"本周期金额,单位：分"`                                                                                                                                             // 本周期金额,单位：分
	Status               int         `json:"status"               description:"订阅单状态，0-Init | 1-Create｜2-Finished｜3-Cancelled"`                                                                                                         // 订阅单状态，0-Init | 1-Create｜2-Finished｜3-Cancelled
	ProrationAmount      int64       `json:"prorationAmount"      description:"下周期金额,单位：分"`                                                                                                                                             // 下周期金额,单位：分
	UpdateAmount         int64       `json:"updateAmount"         description:"升级到金额,单位：分"`                                                                                                                                             // 升级到金额,单位：分
	Currency             string      `json:"currency"             description:"货币"`                                                                                                                                                     // 货币
	UpdateCurrency       string      `json:"updateCurrency"       description:"升级到货币"`                                                                                                                                                  // 升级到货币
	PlanId               int64       `json:"planId"               description:"计划ID"`                                                                                                                                                   // 计划ID
	UpdatePlanId         int64       `json:"updatePlanId"         description:"升级到计划ID"`                                                                                                                                                // 升级到计划ID
	Quantity             int64       `json:"quantity"             description:"quantity"`                                                                                                                                               // quantity
	UpdateQuantity       int64       `json:"updateQuantity"       description:"升级到quantity"`                                                                                                                                            // 升级到quantity
	AddonData            string      `json:"addonData"            description:"plan addon json data"`                                                                                                                                   // plan addon json data
	UpdateAddonData      string      `json:"updateAddonData"      description:"升级到plan addon json data"`                                                                                                                                // 升级到plan addon json data
	ChannelId            int64       `json:"channelId"            description:"支付渠道Id"`                                                                                                                                                 // 支付渠道Id
	UserId               int64       `json:"userId"               description:"userId"`                                                                                                                                                 // userId
	GmtModify            *gtime.Time `json:"gmtModify"            description:"修改时间"`                                                                                                                                                   // 修改时间
	IsDeleted            int         `json:"isDeleted"            description:"0-UnDeleted，1-Deleted"`                                                                                                                                  // 0-UnDeleted，1-Deleted
	Paid                 int         `json:"paid"                 description:"是否已支付，0-否，1-是"`                                                                                                                                          // 是否已支付，0-否，1-是
	Link                 string      `json:"link"                 description:"支付链接"`                                                                                                                                                   // 支付链接
	ChannelStatus        string      `json:"channelStatus"        description:"渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get"` // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	MerchantUserId       int64       `json:"merchantUserId"       description:"merchant_user_id"`                                                                                                                                       // merchant_user_id
	Data                 string      `json:"data"                 description:"渠道额外参数，JSON格式"`                                                                                                                                          // 渠道额外参数，JSON格式
	ResponseData         string      `json:"responseData"         description:"渠道返回参数，JSON格式"`                                                                                                                                          // 渠道返回参数，JSON格式
	EffectImmediate      int         `json:"effectImmediate"      description:"是否马上生效，0-否，1-是"`                                                                                                                                         // 是否马上生效，0-否，1-是
	EffectTime           int64       `json:"effectTime"           description:"effect_immediate=0, 预计生效时间 unit_time"`                                                                                                                   // effect_immediate=0, 预计生效时间 unit_time
	AdminNote            string      `json:"adminNote"            description:"Admin 修改备注"`                                                                                                                                             // Admin 修改备注
	ProrationDate        int64       `json:"prorationDate"        description:"merchant_user_id"`                                                                                                                                       // merchant_user_id
}
