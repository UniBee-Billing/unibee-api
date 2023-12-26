// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Subscription is the golang structure of table subscription for DAO operations like Where/Data.
type Subscription struct {
	g.Meta                `orm:"table:subscription, do:true"`
	Id                    interface{} //
	GmtCreate             *gtime.Time // 创建时间
	GmtModify             *gtime.Time // 修改时间
	CompanyId             interface{} // 公司ID
	MerchantId            interface{} // 商户Id
	PlanId                interface{} // 计划ID
	ChannelId             interface{} // 支付渠道Id
	UserId                interface{} // userId
	Quantity              interface{} // quantity
	SubscriptionId        interface{} // 内部订阅id
	ChannelSubscriptionId interface{} // 支付渠道订阅id
	Data                  interface{} // 渠道额外参数，JSON格式
	ResponseData          interface{} // 渠道返回参数，JSON格式
	IsDeleted             interface{} //
	Status                interface{} // 订阅单状态，0-Init | 1-Create｜2-Active｜3-Inactive
	ChannelUserId         interface{} // 渠道用户 Id
	CustomerName          interface{} // customer_name
	CustomerEmail         interface{} // customer_email
	Link                  interface{} //
}
