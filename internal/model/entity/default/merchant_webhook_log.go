// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantWebhookLog is the golang structure for table merchant_webhook_log.
type MerchantWebhookLog struct {
	Id             uint64      `json:"id"             description:"id"`              // id
	MerchantId     uint64      `json:"merchantId"     description:"webhook url"`     // webhook url
	EndpointId     int64       `json:"endpointId"     description:""`                //
	ReconsumeCount int         `json:"reconsumeCount" description:""`                //
	WebhookUrl     string      `json:"webhookUrl"     description:"webhook url"`     // webhook url
	WebhookEvent   string      `json:"webhookEvent"   description:"webhook_event"`   // webhook_event
	RequestId      string      `json:"requestId"      description:"request_id"`      // request_id
	Body           string      `json:"body"           description:"body(json)"`      // body(json)
	Response       string      `json:"response"       description:"response"`        // response
	Mamo           string      `json:"mamo"           description:"mamo"`            // mamo
	GmtCreate      *gtime.Time `json:"gmtCreate"      description:"create time"`     // create time
	GmtModify      *gtime.Time `json:"gmtModify"      description:"update time"`     // update time
	CreateTime     int64       `json:"createTime"     description:"create utc time"` // create utc time
}
