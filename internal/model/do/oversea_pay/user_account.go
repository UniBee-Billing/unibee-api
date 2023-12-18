// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAccount is the golang structure of table user_account for DAO operations like Where/Data.
type UserAccount struct {
	g.Meta     `orm:"table:user_account, do:true"`
	Id         interface{} // userId
	GmtCreate  *gtime.Time // 创建时间
	GmtModify  *gtime.Time // 修改时间
	IsDeleted  interface{} // 逻辑删除
	Password   interface{} // 密码，加密存储
	UserName   interface{} // 用户名
	Mobile     interface{} // 手机号
	Email      interface{} // 邮箱
	Gender     interface{} // 性别
	AvatarUrl  interface{} // 头像url
	ReMark     interface{} // 备注
	IsSpecial  interface{} // 是否是特殊账号（0.否，1.是）
	Birthday   interface{} // 生日
	Profession interface{} // 职业
	School     interface{} // 学校
	Custom     interface{} // 其他
	NearTime   *gtime.Time // 最近登录时间
	Wid        interface{} // 盟有wid
	IsRisk     interface{} // 风控：0.低风险，1.中风险，2.高风险
	Channel    interface{} // 渠道
}
