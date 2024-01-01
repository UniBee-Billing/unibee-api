// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAccount is the golang structure for table user_account.
type UserAccount struct {
	Id         uint64      `json:"id"         ` // userId
	GmtCreate  *gtime.Time `json:"gmtCreate"  ` // 创建时间
	GmtModify  *gtime.Time `json:"gmtModify"  ` // 修改时间
	IsDeleted  int         `json:"isDeleted"  ` // 逻辑删除
	Password   string      `json:"password"   ` // 密码，加密存储
	UserName   string      `json:"userName"   ` // 用户名
	Mobile     string      `json:"mobile"     ` // 手机号
	Email      string      `json:"email"      ` // 邮箱
	Gender     string      `json:"gender"     ` // 性别
	AvatarUrl  string      `json:"avatarUrl"  ` // 头像url
	ReMark     string      `json:"reMark"     ` // 备注
	IsSpecial  int         `json:"isSpecial"  ` // 是否是特殊账号（0.否，1.是）
	Birthday   string      `json:"birthday"   ` // 生日
	Profession string      `json:"profession" ` // 职业
	School     string      `json:"school"     ` // 学校
	Custom     string      `json:"custom"     ` // 其他
	NearTime   *gtime.Time `json:"nearTime"   ` // 最近登录时间
	Wid        string      `json:"wid"        ` // 盟有wid
	IsRisk     int         `json:"isRisk"     ` // 风控：0.低风险，1.中风险，2.高风险
	Channel    string      `json:"channel"    ` // 渠道
	Version    int         `json:"version"    ` // 版本
	Phone      string      `json:"phone"      ` //
	Address    string      `json:"address"    ` //
	FirstName  string      `json:"firstName"  ` //
	LastName   string      `json:"lastName"   ` //
}
