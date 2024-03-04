// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantWebhookMessage is the golang structure for table merchant_webhook_message.
type MerchantWebhookMessage struct {
	Id              uint64      `json:"id"              description:"id"`                            // id
	MerchantId      uint64      `json:"merchantId"      description:"merchantId"`                    // merchantId
	WebhookEvent    string      `json:"webhookEvent"    description:"webhook_event"`                 // webhook_event
	Data            string      `json:"data"            description:"data(json)"`                    // data(json)
	WebsocketStatus int         `json:"websocketStatus" description:"status  10-pending，20-success"` // status  10-pending，20-success
	GmtCreate       *gtime.Time `json:"gmtCreate"       description:"create time"`                   // create time
	GmtModify       *gtime.Time `json:"gmtModify"       description:"update time"`                   // update time
	CreateTime      int64       `json:"createTime"      description:"create utc time"`               // create utc time
}
