// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantWebhookMessage is the golang structure of table merchant_webhook_message for DAO operations like Where/Data.
type MerchantWebhookMessage struct {
	g.Meta          `orm:"table:merchant_webhook_message, do:true"`
	Id              interface{} // id
	MerchantId      interface{} // merchantId
	WebhookEvent    interface{} // webhook_event
	Data            interface{} // data(json)
	WebsocketStatus interface{} // status  10-pending，20-success
	GmtCreate       *gtime.Time // create time
	GmtModify       *gtime.Time // update time
	CreateTime      interface{} // create utc time
}
