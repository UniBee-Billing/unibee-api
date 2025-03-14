// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionTimeline is the golang structure for table subscription_timeline.
type SubscriptionTimeline struct {
	Id              uint64      `json:"id"              description:""`                        //
	MerchantId      uint64      `json:"merchantId"      description:"merchant id"`             // merchant id
	UserId          uint64      `json:"userId"          description:"userId"`                  // userId
	SubscriptionId  string      `json:"subscriptionId"  description:"subscription id"`         // subscription id
	Status          int         `json:"status"          description:"1-processing,2-finish"`   // 1-processing,2-finish
	PeriodStart     int64       `json:"periodStart"     description:"period_start"`            // period_start
	PeriodEnd       int64       `json:"periodEnd"       description:"period_end"`              // period_end
	PeriodStartTime *gtime.Time `json:"periodStartTime" description:"period start (datetime)"` // period start (datetime)
	PeriodEndTime   *gtime.Time `json:"periodEndTime"   description:"period end (datatime)"`   // period end (datatime)
	GmtCreate       *gtime.Time `json:"gmtCreate"       description:"create time"`             // create time
	GmtModify       *gtime.Time `json:"gmtModify"       description:"update time"`             // update time
	InvoiceId       string      `json:"invoiceId"       description:"invoice id"`              // invoice id
	UniqueId        string      `json:"uniqueId"        description:"unique id"`               // unique id
	Currency        string      `json:"currency"        description:"currency"`                // currency
	PlanId          uint64      `json:"planId"          description:"PlanId"`                  // PlanId
	Quantity        int64       `json:"quantity"        description:"quantity"`                // quantity
	AddonData       string      `json:"addonData"       description:"plan addon json data"`    // plan addon json data
	GatewayId       uint64      `json:"gatewayId"       description:"gateway_id"`              // gateway_id
	IsDeleted       int         `json:"isDeleted"       description:"0-UnDeleted，1-Deleted"`   // 0-UnDeleted，1-Deleted
	UniqueKey       string      `json:"uniqueKey"       description:"unique key (deperated)"`  // unique key (deperated)
	CreateTime      int64       `json:"createTime"      description:"create utc time"`         // create utc time
}
