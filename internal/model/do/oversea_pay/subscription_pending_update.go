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
	Id                   interface{} // id
	MerchantId           interface{} // merchant id
	SubscriptionId       interface{} // subscription id
	UpdateSubscriptionId interface{} // pending update unique id
	GatewayUpdateId      interface{} // gateway update payment id assosiate to this update, use payment.paymentId
	GmtCreate            *gtime.Time // create time
	Amount               interface{} // amount of this period, cent
	Status               interface{} // status，0-Init | 1-Create｜2-Finished｜3-Cancelled
	ProrationAmount      interface{} // proration amount of this pending update , cent
	UpdateAmount         interface{} // the amount after update
	Currency             interface{} // currency of this period
	UpdateCurrency       interface{} // the currency after update
	PlanId               interface{} // the plan id of this period
	UpdatePlanId         interface{} // the plan id after update
	Quantity             interface{} // quantity of this period
	UpdateQuantity       interface{} // quantity after update
	AddonData            interface{} // plan addon data (json) of this period
	UpdateAddonData      interface{} // plan addon data (json) after update
	GatewayId            interface{} // gateway_id
	UserId               interface{} // userId
	GmtModify            *gtime.Time // update time
	IsDeleted            interface{} // 0-UnDeleted，1-Deleted
	Paid                 interface{} // paid，0-no，1-yes
	Link                 interface{} // payment link
	GatewayStatus        interface{} // gateway status
	MerchantUserId       interface{} // merchant_user_id
	Data                 interface{} //
	ResponseData         interface{} //
	EffectImmediate      interface{} // effect immediate，0-no，1-yes
	EffectTime           interface{} // effect_immediate=0, effect time, utc_time
	Note                 interface{} // note
	ProrationDate        interface{} // merchant_user_id
	CreateAt             interface{} // create utc time
}
