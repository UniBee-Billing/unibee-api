// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Timeline is the golang structure for table timeline.
type Timeline struct {
	Id              int64       `json:"id"              ` // 主键id
	UserId          int64       `json:"userId"          ` // user_id
	MerchantUserId  int64       `json:"merchantUserId"  ` // merchant_user_id
	OpenApiId       int64       `json:"openApiId"       ` // 使用的开放平台配置Id
	TerminalIp      string      `json:"terminalIp"      ` // terminal_ip
	BizType         int         `json:"bizType"         ` // biz_type=1，Payment表
	BizId           string      `json:"bizId"           ` // biz_type=1，pay；
	Fee             int64       `json:"fee"             ` // 金额（分）
	EventType       int         `json:"eventType"       ` // 0-未知
	Event           string      `json:"event"           ` // 事件
	RelativeTradeNo string      `json:"relativeTradeNo" ` // 关联单号
	UniqueNo        string      `json:"uniqueNo"        ` // 唯一键
	GmtCreate       *gtime.Time `json:"gmtCreate"       ` // 创建时间
	GmtModify       *gtime.Time `json:"gmtModify"       ` // 更新时间
	Message         string      `json:"message"         ` // message
}
