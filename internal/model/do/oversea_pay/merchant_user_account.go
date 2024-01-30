// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantUserAccount is the golang structure of table merchant_user_account for DAO operations like Where/Data.
type MerchantUserAccount struct {
	g.Meta     `orm:"table:merchant_user_account, do:true"`
	Id         interface{} // userId
	GmtCreate  *gtime.Time // 创建时间
	GmtModify  *gtime.Time // 修改时间
	MerchantId interface{} // 用户ID
	IsDeleted  interface{} // 0-UnDeleted，1-Deleted
	Password   interface{} // 密码，加密存储
	UserName   interface{} // 用户名
	Mobile     interface{} // 手机号
	Email      interface{} // 邮箱
	FirstName  interface{} //
	LastName   interface{} //
}
