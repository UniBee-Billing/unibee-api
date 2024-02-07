// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentTimeline is the golang structure for table payment_timeline.
type PaymentTimeline struct {
	Id             uint64      `json:"id"             description:""`                                //
	MerchantId     int64       `json:"merchantId"     description:"merchant id"`                     // merchant id
	UserId         int64       `json:"userId"         description:"userId"`                          // userId
	SubscriptionId string      `json:"subscriptionId" description:"subscription id"`                 // subscription id
	InvoiceId      string      `json:"invoiceId"      description:"invoice id"`                      // invoice id
	UniqueId       string      `json:"uniqueId"       description:"unique id"`                       // unique id
	Currency       string      `json:"currency"       description:"currency"`                        // currency
	TotalAmount    int64       `json:"totalAmount"    description:"total amount"`                    // total amount
	GatewayId      int64       `json:"gatewayId"      description:"gateway id"`                      // gateway id
	GmtCreate      *gtime.Time `json:"gmtCreate"      description:"create time"`                     // create time
	GmtModify      *gtime.Time `json:"gmtModify"      description:"update time"`                     // update time
	IsDeleted      int         `json:"isDeleted"      description:"0-UnDeleted，1-Deleted"`           // 0-UnDeleted，1-Deleted
	PaymentId      string      `json:"paymentId"      description:"PaymentId"`                       // PaymentId
	Status         int         `json:"status"         description:"0-pending, 1-success, 2-failure"` // 0-pending, 1-success, 2-failure
	TimelineType   int         `json:"timelineType"   description:"0-pay, 1-refund"`                 // 0-pay, 1-refund
	CreateTime     int64       `json:"createTime"     description:"create utc time"`                 // create utc time
}
