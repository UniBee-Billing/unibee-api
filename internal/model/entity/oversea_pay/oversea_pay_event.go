// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OverseaPayEvent is the golang structure for table oversea_pay_event.
type OverseaPayEvent struct {
	Id              int64       `json:"id"              ` // 主键id
	BizType         int         `json:"bizType"         ` // biz_type=0，oversea_pay表
	BizId           int64       `json:"bizId"           ` // biz_type=0，oversea_pay表Id；
	Fee             int64       `json:"fee"             ` // 金额（分）
	EventType       int         `json:"eventType"       ` // 0-未知
	Event           string      `json:"event"           ` // 事件
	RelativeTradeNo string      `json:"relativeTradeNo" ` // 关联单号
	UniqueNo        string      `json:"uniqueNo"        ` // 唯一键
	GmtCreate       *gtime.Time `json:"gmtCreate"       ` // 创建时间
	GmtModify       *gtime.Time `json:"gmtModify"       ` // 更新时间
	OpenApiId       int64       `json:"openApiId"       ` // 使用的开放平台配置Id
	TerminalIp      string      `json:"terminalIp"      ` // 实时交易终端IP
	Message         string      `json:"message"         ` // message
}
