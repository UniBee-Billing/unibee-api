// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantWebhook is the golang structure for table merchant_webhook.
type MerchantWebhook struct {
	Id            uint64      `json:"id"            description:"id"`                       // id
	MerchantId    int64       `json:"merchantId"    description:"webhook url"`              // webhook url
	WebhookUrl    string      `json:"webhookUrl"    description:"webhook url"`              // webhook url
	WebhookEvents string      `json:"webhookEvents" description:"webhook_events,split dot"` // webhook_events,split dot
	GmtCreate     *gtime.Time `json:"gmtCreate"     description:"create time"`              // create time
	GmtModify     *gtime.Time `json:"gmtModify"     description:"update time"`              // update time
	CreateTime    int64       `json:"createTime"    description:"create utc time"`          // create utc time
}
