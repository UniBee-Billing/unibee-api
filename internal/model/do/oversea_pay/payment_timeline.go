// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentTimeline is the golang structure of table payment_timeline for DAO operations like Where/Data.
type PaymentTimeline struct {
	g.Meta          `orm:"table:payment_timeline, do:true"`
	Id              interface{} // 主键id
	UserId          interface{} // user_id
	MerchantUserId  interface{} // merchant_user_id
	OpenApiId       interface{} // 使用的开放平台配置Id
	TerminalIp      interface{} // terminal_ip
	BizType         interface{} // biz_type=1，Payment表
	BizId           interface{} // biz_type=1，pay；
	Fee             interface{} // 金额（分）
	EventType       interface{} // 0-未知
	Event           interface{} // 事件
	RelativeTradeNo interface{} // 关联单号
	UniqueNo        interface{} // 唯一键
	GmtCreate       *gtime.Time // 创建时间
	GmtModify       *gtime.Time // 更新时间
	Message         interface{} // message
}
