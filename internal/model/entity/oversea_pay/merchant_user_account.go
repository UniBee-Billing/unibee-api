// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantUserAccount is the golang structure for table merchant_user_account.
type MerchantUserAccount struct {
	Id         uint64      `json:"id"         ` // userId
	GmtCreate  *gtime.Time `json:"gmtCreate"  ` // 创建时间
	GmtModify  *gtime.Time `json:"gmtModify"  ` // 修改时间
	MerchantId int64       `json:"merchantId" ` // 用户ID
	IsDeleted  int         `json:"isDeleted"  ` // 逻辑删除
	Password   string      `json:"password"   ` // 密码，加密存储
	UserName   string      `json:"userName"   ` // 用户名
	Mobile     string      `json:"mobile"     ` // 手机号
	Email      string      `json:"email"      ` // 邮箱
	FirstName  string      `json:"firstName"  ` //
	LastName   string      `json:"lastName"   ` //
}
