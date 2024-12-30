// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantDiscountCode is the golang structure for table merchant_discount_code.
type MerchantDiscountCode struct {
	Id                 uint64      `json:"id"                 description:"ID"`                                                                                                                       // ID
	MerchantId         uint64      `json:"merchantId"         description:"merchantId"`                                                                                                               // merchantId
	Name               string      `json:"name"               description:"name"`                                                                                                                     // name
	Code               string      `json:"code"               description:"code"`                                                                                                                     // code
	Status             int         `json:"status"             description:"status, 1-editable, 2-active, 3-deactive, 4-expire"`                                                                       // status, 1-editable, 2-active, 3-deactive, 4-expire
	BillingType        int         `json:"billingType"        description:"billing_type, 1-one-time, 2-recurring"`                                                                                    // billing_type, 1-one-time, 2-recurring
	DiscountType       int         `json:"discountType"       description:"discount_type, 1-percentage, 2-fixed_amount"`                                                                              // discount_type, 1-percentage, 2-fixed_amount
	DiscountAmount     int64       `json:"discountAmount"     description:"amount of discount, available when discount_type is fixed_amount"`                                                         // amount of discount, available when discount_type is fixed_amount
	DiscountPercentage int64       `json:"discountPercentage" description:"percentage of discount, 100=1%, available when discount_type is percentage"`                                               // percentage of discount, 100=1%, available when discount_type is percentage
	Currency           string      `json:"currency"           description:"currency of discount, available when discount_type is fixed_amount"`                                                       // currency of discount, available when discount_type is fixed_amount
	SubscriptionLimit  int         `json:"subscriptionLimit"  description:"the limit of every subscription apply, 0-unlimited"`                                                                       // the limit of every subscription apply, 0-unlimited
	StartTime          int64       `json:"startTime"          description:"start of discount available utc time"`                                                                                     // start of discount available utc time
	EndTime            int64       `json:"endTime"            description:"end of discount available utc time, 0-invalid"`                                                                            // end of discount available utc time, 0-invalid
	GmtCreate          *gtime.Time `json:"gmtCreate"          description:"create time"`                                                                                                              // create time
	GmtModify          *gtime.Time `json:"gmtModify"          description:"update time"`                                                                                                              // update time
	IsDeleted          int         `json:"isDeleted"          description:"0-UnDeleted，1-Deleted"`                                                                                                    // 0-UnDeleted，1-Deleted
	CreateTime         int64       `json:"createTime"         description:"create utc time"`                                                                                                          // create utc time
	CycleLimit         int         `json:"cycleLimit"         description:"the count limitation of subscription cycle , 0-no limit"`                                                                  // the count limitation of subscription cycle , 0-no limit
	MetaData           string      `json:"metaData"           description:"meta_data(json)"`                                                                                                          // meta_data(json)
	Type               int         `json:"type"               description:"type, 1-external discount code"`                                                                                           // type, 1-external discount code
	PlanIds            string      `json:"planIds"            description:"Ids of plan which discount code can effect, default effect all plans if not set"`                                          // Ids of plan which discount code can effect, default effect all plans if not set
	Quantity           int64       `json:"quantity"           description:"quantity of code"`                                                                                                         // quantity of code
	Advance            int         `json:"advance"            description:"AdvanceConfig,  0-false,1-true, will enable all advance config if set 1"`                                                  // AdvanceConfig,  0-false,1-true, will enable all advance config if set 1
	UserLimit          int         `json:"userLimit"          description:"AdvanceConfig, The limit of every customer can apply, the recurring apply not involved, 0-unlimited\""`                    // AdvanceConfig, The limit of every customer can apply, the recurring apply not involved, 0-unlimited"
	UserScope          int         `json:"userScope"          description:"AdvanceConfig, Apply user scope,0-for all, 1-for only new user, 2-for only renewals, renewals is upgrade&downgrade&renew"` // AdvanceConfig, Apply user scope,0-for all, 1-for only new user, 2-for only renewals, renewals is upgrade&downgrade&renew
	UpgradeOnly        int         `json:"upgradeOnly"        description:"AdvanceConfig, 0-false,1-true, will forbid for all except upgrade action if set 1"`                                        // AdvanceConfig, 0-false,1-true, will forbid for all except upgrade action if set 1
	UpgradeLongerOnly  int         `json:"upgradeLongerOnly"  description:"AdvanceConfig, 0-false,1-true, will forbid for all except upgrade to longer plan if set 1"`                                // AdvanceConfig, 0-false,1-true, will forbid for all except upgrade to longer plan if set 1
}
