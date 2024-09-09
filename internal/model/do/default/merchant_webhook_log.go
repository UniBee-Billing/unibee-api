// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantWebhookLog is the golang structure of table merchant_webhook_log for DAO operations like Where/Data.
type MerchantWebhookLog struct {
	g.Meta         `orm:"table:merchant_webhook_log, do:true"`
	Id             interface{} // id
	MerchantId     interface{} // webhook url
	EndpointId     interface{} //
	ReconsumeCount interface{} //
	WebhookUrl     interface{} // webhook url
	WebhookEvent   interface{} // webhook_event
	RequestId      interface{} // request_id
	Body           interface{} // body(json)
	Response       interface{} // response
	Mamo           interface{} // mamo
	GmtCreate      *gtime.Time // create time
	GmtModify      *gtime.Time // update time
	CreateTime     interface{} // create utc time
	WebhookEventId interface{} // webhook_event_id
}
