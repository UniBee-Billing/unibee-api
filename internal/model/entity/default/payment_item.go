// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentItem is the golang structure for table payment_item.
type PaymentItem struct {
	Id             uint64      `json:"id"             description:""`                                           //
	BizType        int         `json:"bizType"        description:"biz_type 1-onetime payment, 3-subscription"` // biz_type 1-onetime payment, 3-subscription
	Status         int         `json:"status"         description:"0-pending, 1-success, 2-failure"`            // 0-pending, 1-success, 2-failure
	MerchantId     uint64      `json:"merchantId"     description:"merchant id"`                                // merchant id
	UserId         uint64      `json:"userId"         description:"userId"`                                     // userId
	SubscriptionId string      `json:"subscriptionId" description:"subscription id"`                            // subscription id
	InvoiceId      string      `json:"invoiceId"      description:"invoice id"`                                 // invoice id
	UniqueId       string      `json:"uniqueId"       description:"unique id"`                                  // unique id
	Currency       string      `json:"currency"       description:"currency"`                                   // currency
	Amount         int64       `json:"amount"         description:"amount"`                                     // amount
	UnitAmount     int64       `json:"unitAmount"     description:"unit_amount"`                                // unit_amount
	Quantity       int64       `json:"quantity"       description:"quantity"`                                   // quantity
	GatewayId      uint64      `json:"gatewayId"      description:"gateway id"`                                 // gateway id
	GmtCreate      *gtime.Time `json:"gmtCreate"      description:"create time"`                                // create time
	GmtModify      *gtime.Time `json:"gmtModify"      description:"update time"`                                // update time
	IsDeleted      int         `json:"isDeleted"      description:"0-UnDeleted，1-Deleted"`                      // 0-UnDeleted，1-Deleted
	PaymentId      string      `json:"paymentId"      description:"PaymentId"`                                  // PaymentId
	CreateTime     int64       `json:"createTime"     description:"create utc time"`                            // create utc time
	Description    string      `json:"description"    description:"description"`                                // description
	Name           string      `json:"name"           description:"name"`                                       // name
}
