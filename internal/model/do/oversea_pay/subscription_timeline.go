// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionTimeline is the golang structure of table subscription_timeline for DAO operations like Where/Data.
type SubscriptionTimeline struct {
	g.Meta          `orm:"table:subscription_timeline, do:true"`
	Id              interface{} //
	MerchantId      interface{} // merchant id
	UserId          interface{} // userId
	SubscriptionId  interface{} // subscription id
	PeriodStart     interface{} // period_start
	PeriodEnd       interface{} // period_end
	PeriodStartTime *gtime.Time // period start (datetime)
	PeriodEndTime   *gtime.Time // period end (datatime)
	GmtCreate       *gtime.Time // create time
	InvoiceId       interface{} // invoice id
	UniqueId        interface{} // unique id
	Currency        interface{} // currency
	PlanId          interface{} // PlanId
	Quantity        interface{} // quantity
	AddonData       interface{} // plan addon json data
	ChannelId       interface{} // channel_id
	GmtModify       *gtime.Time // update time
	IsDeleted       interface{} // 0-UnDeleted，1-Deleted
	UniqueKey       interface{} // unique key (deperated)
}
