// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentEvent is the golang structure of table payment_event for DAO operations like Where/Data.
type PaymentEvent struct {
	g.Meta          `orm:"table:payment_event, do:true"`
	Id              interface{} // id
	UserId          interface{} // user_id
	MerchantUserId  interface{} // merchant_user_id
	OpenApiId       interface{} // open api id
	TerminalIp      interface{} // terminal_ip
	BizType         interface{} // biz_type=1，Payment表
	BizId           interface{} // biz_type=1，pay；
	Fee             interface{} // amount, cent
	EventType       interface{} // 0-unknown
	Event           interface{} // event
	RelativeTradeNo interface{} // relative trade no
	UniqueNo        interface{} // unique no
	GmtCreate       *gtime.Time // create time
	GmtModify       *gtime.Time // update time
	Message         interface{} // message
	CreateAt        interface{} // create utc time
}
