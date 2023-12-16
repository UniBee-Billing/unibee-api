// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OverseaPayEvent is the golang structure of table oversea_pay_event for DAO operations like Where/Data.
type OverseaPayEvent struct {
	g.Meta          `orm:"table:oversea_pay_event, do:true"`
	Id              interface{} // 主键id
	BizType         interface{} // biz_type=0，oversea_pay表
	BizId           interface{} // biz_type=0，oversea_pay表Id；
	Fee             interface{} // 金额（分）
	EventType       interface{} // 0-未知
	Event           interface{} // 事件
	RelativeTradeNo interface{} // 关联单号
	UniqueNo        interface{} // 唯一键
	GmtCreate       *gtime.Time // 创建时间
	GmtModify       *gtime.Time // 更新时间
	OpenApiId       interface{} // 使用的开放平台配置Id
	TerminalIp      interface{} // 实时交易终端IP
	Message         interface{} // message
}
