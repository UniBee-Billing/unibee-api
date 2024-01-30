// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantInfo is the golang structure of table merchant_info for DAO operations like Where/Data.
type MerchantInfo struct {
	g.Meta      `orm:"table:merchant_info, do:true"`
	Id          interface{} // 用户的ID
	CompanyId   interface{} // 公司ID
	UserId      interface{} // 用户ID
	Type        interface{} // 类型，0-个人，1-企业
	CompanyName interface{} // 公司名称
	Email       interface{} // email
	BusinessNum interface{} // 税号
	Name        interface{} // 个人或法人姓名
	Idcard      interface{} // 个人或法人身份证号
	Location    interface{} // 省市区地址
	Address     interface{} // 详细地址
	GmtCreate   *gtime.Time // 创建时间
	GmtModify   *gtime.Time // 修改时间
	IsDeleted   interface{} // 0-UnDeleted，1-Deleted
	CompanyLogo interface{} // 账号头像
	HomeUrl     interface{} //
	FirstName   interface{} //
	LastName    interface{} //
	Phone       interface{} //
}
