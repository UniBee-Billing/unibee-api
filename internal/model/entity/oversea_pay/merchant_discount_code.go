// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantDiscountCode is the golang structure for table merchant_discount_code.
type MerchantDiscountCode struct {
	Id                int64       `json:"id"                description:"ID"`                                                                // ID
	MerchantId        uint64      `json:"merchantId"        description:"merchantId"`                                                        // merchantId
	Name              string      `json:"name"              description:"name"`                                                              // name
	Code              string      `json:"code"              description:"code"`                                                              // code
	Status            int         `json:"status"            description:"status, 1-editable, 2-active, 3-deactive, 4-expire"`                // status, 1-editable, 2-active, 3-deactive, 4-expire
	BillingType       int         `json:"billingType"       description:"billing_type, 1-one-time, 2-recurring"`                             // billing_type, 1-one-time, 2-recurring
	DiscountType      int         `json:"discountType"      description:"discount_type, 1-percentage, 2-fixed_amount"`                       // discount_type, 1-percentage, 2-fixed_amount
	Amount            int64       `json:"amount"            description:"amount of discount, avalible when discount_type is fixed_amount"`   // amount of discount, avalible when discount_type is fixed_amount
	Currency          string      `json:"currency"          description:"currency of discount, avalible when discount_type is fixed_amount"` // currency of discount, avalible when discount_type is fixed_amount
	UserLimit         int         `json:"userLimit"         description:"the limit of every user apply, 0-unlimit"`                          // the limit of every user apply, 0-unlimit
	SubscriptionLimit int         `json:"subscriptionLimit" description:"the limit of every subscription apply, 0-unlimit"`                  // the limit of every subscription apply, 0-unlimit
	StartTime         int64       `json:"startTime"         description:"start of discount avalible utc time"`                               // start of discount avalible utc time
	EndTime           int64       `json:"endTime"           description:"end of discount avalible utc time"`                                 // end of discount avalible utc time
	GmtCreate         *gtime.Time `json:"gmtCreate"         description:"create time"`                                                       // create time
	GmtModify         *gtime.Time `json:"gmtModify"         description:"update time"`                                                       // update time
	IsDeleted         int         `json:"isDeleted"         description:"0-UnDeleted，1-Deleted"`                                             // 0-UnDeleted，1-Deleted
	CreateTime        int64       `json:"createTime"        description:"create utc time"`                                                   // create utc time
}
