// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionPendingUpdate is the golang structure for table subscription_pending_update.
type SubscriptionPendingUpdate struct {
	Id                   uint64      `json:"id"                   description:"id"`                                                                        // id
	MerchantId           int64       `json:"merchantId"           description:"merchant id"`                                                               // merchant id
	SubscriptionId       string      `json:"subscriptionId"       description:"subscription id"`                                                           // subscription id
	UpdateSubscriptionId string      `json:"updateSubscriptionId" description:"pending update unique id"`                                                  // pending update unique id
	ChannelUpdateId      string      `json:"channelUpdateId"      description:"channel update payment id assosiate to this update, use payment.paymentId"` // channel update payment id assosiate to this update, use payment.paymentId
	GmtCreate            *gtime.Time `json:"gmtCreate"            description:"create time"`                                                               // create time
	Amount               int64       `json:"amount"               description:"amount of this period, cent"`                                               // amount of this period, cent
	Status               int         `json:"status"               description:"status，0-Init | 1-Create｜2-Finished｜3-Cancelled"`                           // status，0-Init | 1-Create｜2-Finished｜3-Cancelled
	ProrationAmount      int64       `json:"prorationAmount"      description:"proration amount of this pending update , cent"`                            // proration amount of this pending update , cent
	UpdateAmount         int64       `json:"updateAmount"         description:"the amount after update"`                                                   // the amount after update
	Currency             string      `json:"currency"             description:"currency of this period"`                                                   // currency of this period
	UpdateCurrency       string      `json:"updateCurrency"       description:"the currency after update"`                                                 // the currency after update
	PlanId               int64       `json:"planId"               description:"the plan id of this period"`                                                // the plan id of this period
	UpdatePlanId         int64       `json:"updatePlanId"         description:"the plan id after update"`                                                  // the plan id after update
	Quantity             int64       `json:"quantity"             description:"quantity of this period"`                                                   // quantity of this period
	UpdateQuantity       int64       `json:"updateQuantity"       description:"quantity after update"`                                                     // quantity after update
	AddonData            string      `json:"addonData"            description:"plan addon data (json) of this period"`                                     // plan addon data (json) of this period
	UpdateAddonData      string      `json:"updateAddonData"      description:"plan addon data (json) after update"`                                       // plan addon data (json) after update
	ChannelId            int64       `json:"channelId"            description:"channel_id"`                                                                // channel_id
	UserId               int64       `json:"userId"               description:"userId"`                                                                    // userId
	GmtModify            *gtime.Time `json:"gmtModify"            description:"update time"`                                                               // update time
	IsDeleted            int         `json:"isDeleted"            description:"0-UnDeleted，1-Deleted"`                                                     // 0-UnDeleted，1-Deleted
	Paid                 int         `json:"paid"                 description:"paid，0-no，1-yes"`                                                           // paid，0-no，1-yes
	Link                 string      `json:"link"                 description:"payment link"`                                                              // payment link
	ChannelStatus        string      `json:"channelStatus"        description:"channel status"`                                                            // channel status
	MerchantUserId       int64       `json:"merchantUserId"       description:"merchant_user_id"`                                                          // merchant_user_id
	Data                 string      `json:"data"                 description:""`                                                                          //
	ResponseData         string      `json:"responseData"         description:""`                                                                          //
	EffectImmediate      int         `json:"effectImmediate"      description:"是否马上生效，0-否，1-是"`                                                            // 是否马上生效，0-否，1-是
	EffectTime           int64       `json:"effectTime"           description:"effect_immediate=0, 预计生效时间 unit_time"`                                      // effect_immediate=0, 预计生效时间 unit_time
	Note                 string      `json:"note"                 description:"note"`                                                                      // note
	ProrationDate        int64       `json:"prorationDate"        description:"merchant_user_id"`                                                          // merchant_user_id
}
