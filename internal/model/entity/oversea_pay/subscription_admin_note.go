// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionAdminNote is the golang structure for table subscription_admin_note.
type SubscriptionAdminNote struct {
	Id             uint64      `json:"id"             description:""`                 //
	GmtCreate      *gtime.Time `json:"gmtCreate"      description:"create_time"`      // create_time
	GmtModify      *gtime.Time `json:"gmtModify"      description:"modify_time"`      // modify_time
	SubscriptionId string      `json:"subscriptionId" description:"subscription_id"`  // subscription_id
	MerchantUserId int64       `json:"merchantUserId" description:"merchant_user_id"` // merchant_user_id
	Note           string      `json:"note"           description:"note"`             // note
	IsDeleted      int         `json:"isDeleted"      description:""`                 //
}
