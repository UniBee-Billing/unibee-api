// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionOnetimeAddon is the golang structure for table subscription_onetime_addon.
type SubscriptionOnetimeAddon struct {
	Id             uint64      `json:"id"             description:"id"`                                            // id
	GmtCreate      *gtime.Time `json:"gmtCreate"      description:"create_time"`                                   // create_time
	GmtModify      *gtime.Time `json:"gmtModify"      description:"modify_time"`                                   // modify_time
	SubscriptionId string      `json:"subscriptionId" description:"subscription_id"`                               // subscription_id
	AddonId        uint64      `json:"addonId"        description:"onetime addonId"`                               // onetime addonId
	Quantity       int64       `json:"quantity"       description:"quantity"`                                      // quantity
	Status         int         `json:"status"         description:"status, 1-create, 2-paid, 3-cancel, 4-expired"` // status, 1-create, 2-paid, 3-cancel, 4-expired
	IsDeleted      int         `json:"isDeleted"      description:"0-UnDeleted，1-Deleted"`                         // 0-UnDeleted，1-Deleted
	CreateTime     int64       `json:"createTime"     description:"create utc time"`                               // create utc time
	PaymentId      string      `json:"paymentId"      description:"paymentId"`                                     // paymentId
	MetaData       string      `json:"metaData"       description:"meta_data(json)"`                               // meta_data(json)
}
