// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionPendingUpdate is the golang structure for table subscription_pending_update.
type SubscriptionPendingUpdate struct {
	Id               uint64      `json:"id"               description:"id"`                                                                        // id
	MerchantId       uint64      `json:"merchantId"       description:"merchant id"`                                                               // merchant id
	SubscriptionId   string      `json:"subscriptionId"   description:"subscription id"`                                                           // subscription id
	PendingUpdateId  string      `json:"pendingUpdateId"  description:"pending update unique id"`                                                  // pending update unique id
	Name             string      `json:"name"             description:"name"`                                                                      // name
	InvoiceId        string      `json:"invoiceId"        description:"gateway update payment id assosiate to this update, use payment.paymentId"` // gateway update payment id assosiate to this update, use payment.paymentId
	GmtCreate        *gtime.Time `json:"gmtCreate"        description:"create time"`                                                               // create time
	GmtModify        *gtime.Time `json:"gmtModify"        description:"update time"`                                                               // update time
	Amount           int64       `json:"amount"           description:"amount of this period, cent"`                                               // amount of this period, cent
	Status           int         `json:"status"           description:"status，0-Init | 1-Pending｜2-Finished｜3-Cancelled"`                          // status，0-Init | 1-Pending｜2-Finished｜3-Cancelled
	ProrationAmount  int64       `json:"prorationAmount"  description:"proration amount of this pending update , cent"`                            // proration amount of this pending update , cent
	UpdateAmount     int64       `json:"updateAmount"     description:"the amount after update"`                                                   // the amount after update
	Currency         string      `json:"currency"         description:"currency of this period"`                                                   // currency of this period
	UpdateCurrency   string      `json:"updateCurrency"   description:"the currency after update"`                                                 // the currency after update
	PlanId           uint64      `json:"planId"           description:"the plan id of this period"`                                                // the plan id of this period
	UpdatePlanId     uint64      `json:"updatePlanId"     description:"the plan id after update"`                                                  // the plan id after update
	Quantity         int64       `json:"quantity"         description:"quantity of this period"`                                                   // quantity of this period
	UpdateQuantity   int64       `json:"updateQuantity"   description:"quantity after update"`                                                     // quantity after update
	AddonData        string      `json:"addonData"        description:"plan addon data (json) of this period"`                                     // plan addon data (json) of this period
	UpdateAddonData  string      `json:"updateAddonData"  description:"plan addon data (json) after update"`                                       // plan addon data (json) after update
	GatewayId        uint64      `json:"gatewayId"        description:"gateway_id"`                                                                // gateway_id
	UserId           uint64      `json:"userId"           description:"userId"`                                                                    // userId
	IsDeleted        int         `json:"isDeleted"        description:"0-UnDeleted，1-Deleted"`                                                     // 0-UnDeleted，1-Deleted
	Paid             int         `json:"paid"             description:"paid，0-no，1-yes"`                                                           // paid，0-no，1-yes
	Link             string      `json:"link"             description:"payment link"`                                                              // payment link
	GatewayStatus    string      `json:"gatewayStatus"    description:"gateway status"`                                                            // gateway status
	MerchantMemberId int64       `json:"merchantMemberId" description:"merchant_user_id"`                                                          // merchant_user_id
	Data             string      `json:"data"             description:""`                                                                          //
	ResponseData     string      `json:"responseData"     description:""`                                                                          //
	EffectImmediate  int         `json:"effectImmediate"  description:"effect immediate，0-no，1-yes"`                                               // effect immediate，0-no，1-yes
	EffectTime       int64       `json:"effectTime"       description:"effect_immediate=0, effect time, utc_time"`                                 // effect_immediate=0, effect time, utc_time
	Note             string      `json:"note"             description:"note"`                                                                      // note
	ProrationDate    int64       `json:"prorationDate"    description:"merchant_user_id"`                                                          // merchant_user_id
	CreateTime       int64       `json:"createTime"       description:"create utc time"`                                                           // create utc time
	MetaData         string      `json:"metaData"         description:"meta_data(json)"`                                                           // meta_data(json)
	DiscountCode     string      `json:"discountCode"     description:"discount_code"`                                                             // discount_code
	TaxPercentage    int64       `json:"taxPercentage"    description:"taxPercentage，1000 = 10%"`                                                  // taxPercentage，1000 = 10%
}
