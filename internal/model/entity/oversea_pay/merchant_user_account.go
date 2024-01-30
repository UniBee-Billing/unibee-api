// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantUserAccount is the golang structure for table merchant_user_account.
type MerchantUserAccount struct {
	Id         uint64      `json:"id"         description:"userId"`                // userId
	GmtCreate  *gtime.Time `json:"gmtCreate"  description:"创建时间"`                  // 创建时间
	GmtModify  *gtime.Time `json:"gmtModify"  description:"修改时间"`                  // 修改时间
	MerchantId int64       `json:"merchantId" description:"用户ID"`                  // 用户ID
	IsDeleted  int         `json:"isDeleted"  description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	Password   string      `json:"password"   description:"密码，加密存储"`               // 密码，加密存储
	UserName   string      `json:"userName"   description:"用户名"`                   // 用户名
	Mobile     string      `json:"mobile"     description:"手机号"`                   // 手机号
	Email      string      `json:"email"      description:"邮箱"`                    // 邮箱
	FirstName  string      `json:"firstName"  description:""`                      //
	LastName   string      `json:"lastName"   description:""`                      //
}
