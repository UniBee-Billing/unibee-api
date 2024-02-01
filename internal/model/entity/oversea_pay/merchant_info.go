// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantInfo is the golang structure for table merchant_info.
type MerchantInfo struct {
	Id          int64       `json:"id"          description:"用户的ID"`                 // 用户的ID
	CompanyId   int64       `json:"companyId"   description:"公司ID"`                  // 公司ID
	UserId      int64       `json:"userId"      description:"用户ID"`                  // 用户ID
	Type        int         `json:"type"        description:"类型，0-个人，1-企业"`          // 类型，0-个人，1-企业
	CompanyName string      `json:"companyName" description:"公司名称"`                  // 公司名称
	Email       string      `json:"email"       description:"email"`                 // email
	BusinessNum string      `json:"businessNum" description:"税号"`                    // 税号
	Name        string      `json:"name"        description:"个人或法人姓名"`               // 个人或法人姓名
	Idcard      string      `json:"idcard"      description:"个人或法人身份证号"`             // 个人或法人身份证号
	Location    string      `json:"location"    description:"省市区地址"`                 // 省市区地址
	Address     string      `json:"address"     description:"详细地址"`                  // 详细地址
	GmtCreate   *gtime.Time `json:"gmtCreate"   description:"create time"`           // create time
	GmtModify   *gtime.Time `json:"gmtModify"   description:"修改时间"`                  // 修改时间
	IsDeleted   int         `json:"isDeleted"   description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	CompanyLogo string      `json:"companyLogo" description:"账号头像"`                  // 账号头像
	HomeUrl     string      `json:"homeUrl"     description:""`                      //
	FirstName   string      `json:"firstName"   description:""`                      //
	LastName    string      `json:"lastName"    description:""`                      //
	Phone       string      `json:"phone"       description:""`                      //
}
