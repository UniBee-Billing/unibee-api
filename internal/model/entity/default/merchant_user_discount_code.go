// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantUserDiscountCode is the golang structure for table merchant_user_discount_code.
type MerchantUserDiscountCode struct {
	Id             int64       `json:"id"             description:"ID"`                                                           // ID
	MerchantId     uint64      `json:"merchantId"     description:"merchantId"`                                                   // merchantId
	UserId         uint64      `json:"userId"         description:"user_id"`                                                      // user_id
	Code           string      `json:"code"           description:"code"`                                                         // code
	Status         int         `json:"status"         description:"status, 1-normal, 2-rollback"`                                 // status, 1-normal, 2-rollback
	PlanId         string      `json:"planId"         description:"plan_id"`                                                      // plan_id
	SubscriptionId string      `json:"subscriptionId" description:"subscription_id"`                                              // subscription_id
	PaymentId      string      `json:"paymentId"      description:"payment_id"`                                                   // payment_id
	InvoiceId      string      `json:"invoiceId"      description:"invoice_id"`                                                   // invoice_id
	UniqueId       string      `json:"uniqueId"       description:"unique_id"`                                                    // unique_id
	GmtCreate      *gtime.Time `json:"gmtCreate"      description:"create time"`                                                  // create time
	GmtModify      *gtime.Time `json:"gmtModify"      description:"update time"`                                                  // update time
	IsDeleted      int         `json:"isDeleted"      description:"0-UnDeleted，1-Deleted"`                                        // 0-UnDeleted，1-Deleted
	CreateTime     int64       `json:"createTime"     description:"create utc time"`                                              // create utc time
	ApplyAmount    int64       `json:"applyAmount"    description:"apply_amount"`                                                 // apply_amount
	Currency       string      `json:"currency"       description:"currency"`                                                     // currency
	Recurring      int         `json:"recurring"      description:"is recurring apply, 0-no, 1-yes"`                              // is recurring apply, 0-no, 1-yes
	RecurringId    int64       `json:"recurringId"    description:"the first purchase id for the code, using for reucrring code"` // the first purchase id for the code, using for reucrring code
}
