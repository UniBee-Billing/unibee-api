// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentTimeline is the golang structure for table payment_timeline.
type PaymentTimeline struct {
	Id              int64       `json:"id"              description:"主键id"`                // 主键id
	UserId          int64       `json:"userId"          description:"user_id"`             // user_id
	MerchantUserId  int64       `json:"merchantUserId"  description:"merchant_user_id"`    // merchant_user_id
	OpenApiId       int64       `json:"openApiId"       description:"使用的开放平台配置Id"`         // 使用的开放平台配置Id
	TerminalIp      string      `json:"terminalIp"      description:"terminal_ip"`         // terminal_ip
	BizType         int         `json:"bizType"         description:"biz_type=1，Payment表"` // biz_type=1，Payment表
	BizId           string      `json:"bizId"           description:"biz_type=1，pay；"`     // biz_type=1，pay；
	Fee             int64       `json:"fee"             description:"金额（分）"`               // 金额（分）
	EventType       int         `json:"eventType"       description:"0-未知"`                // 0-未知
	Event           string      `json:"event"           description:"事件"`                  // 事件
	RelativeTradeNo string      `json:"relativeTradeNo" description:"关联单号"`                // 关联单号
	UniqueNo        string      `json:"uniqueNo"        description:"唯一键"`                 // 唯一键
	GmtCreate       *gtime.Time `json:"gmtCreate"       description:"创建时间"`                // 创建时间
	GmtModify       *gtime.Time `json:"gmtModify"       description:"更新时间"`                // 更新时间
	Message         string      `json:"message"         description:"message"`             // message
}
