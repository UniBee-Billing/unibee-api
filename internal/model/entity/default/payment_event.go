// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentEvent is the golang structure for table payment_event.
type PaymentEvent struct {
	Id               int64       `json:"id"               description:"id"`                 // id
	UserId           int64       `json:"userId"           description:"user_id"`            // user_id
	MerchantMemberId int64       `json:"merchantMemberId" description:"merchant_user_id"`   // merchant_user_id
	OpenApiId        int64       `json:"openApiId"        description:"open api id"`        // open api id
	TerminalIp       string      `json:"terminalIp"       description:"terminal_ip"`        // terminal_ip
	BizType          int         `json:"bizType"          description:"biz_type=1，Payment"` // biz_type=1，Payment
	BizId            string      `json:"bizId"            description:"biz_type=1，pay；"`    // biz_type=1，pay；
	Fee              int64       `json:"fee"              description:"amount, cent"`       // amount, cent
	EventType        int         `json:"eventType"        description:"0-unknown"`          // 0-unknown
	Event            string      `json:"event"            description:"event"`              // event
	RelativeTradeNo  string      `json:"relativeTradeNo"  description:"relative trade no"`  // relative trade no
	UniqueNo         string      `json:"uniqueNo"         description:"unique no"`          // unique no
	GmtCreate        *gtime.Time `json:"gmtCreate"        description:"create time"`        // create time
	GmtModify        *gtime.Time `json:"gmtModify"        description:"update time"`        // update time
	Message          string      `json:"message"          description:"message"`            // message
	CreateTime       int64       `json:"createTime"       description:"create utc time"`    // create utc time
}
