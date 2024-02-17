// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantWebhook is the golang structure of table merchant_webhook for DAO operations like Where/Data.
type MerchantWebhook struct {
	g.Meta        `orm:"table:merchant_webhook, do:true"`
	Id            interface{} // id
	MerchantId    interface{} // webhook url
	WebhookUrl    interface{} // webhook url
	WebhookEvents interface{} // webhook_events,split dot
	GmtCreate     *gtime.Time // create time
	GmtModify     *gtime.Time // update time
	CreateTime    interface{} // create utc time
}
