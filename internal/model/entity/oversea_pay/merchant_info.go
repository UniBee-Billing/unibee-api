// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantInfo is the golang structure for table merchant_info.
type MerchantInfo struct {
	Id          int64       `json:"id"          ` // 用户的ID
	CompanyId   int64       `json:"companyId"   ` // 公司ID
	UserId      int64       `json:"userId"      ` // 用户ID
	Type        int         `json:"type"        ` // 类型，0-个人，1-企业
	CompanyName string      `json:"companyName" ` // 公司名称
	Email       string      `json:"email"       ` // email
	BusinessNum string      `json:"businessNum" ` // 税号
	Name        string      `json:"name"        ` // 个人或法人姓名
	Idcard      string      `json:"idcard"      ` // 个人或法人身份证号
	Location    string      `json:"location"    ` // 省市区地址
	Address     string      `json:"address"     ` // 详细地址
	GmtCreate   *gtime.Time `json:"gmtCreate"   ` // 创建时间
	GmtModify   *gtime.Time `json:"gmtModify"   ` // 修改时间
	IsDeleted   int         `json:"isDeleted"   ` // 是否删除，0-未删除，1-删除
	CompanyLogo string      `json:"companyLogo" ` // 账号头像
	HomeUrl     string      `json:"homeUrl"     ` //
}
