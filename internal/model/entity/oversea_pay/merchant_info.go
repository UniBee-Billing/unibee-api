// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantInfo is the golang structure for table merchant_info.
type MerchantInfo struct {
	Id                     int64       `json:"id"                     ` // 用户的ID
	CompanyId              int64       `json:"companyId"              ` // 公司ID
	UserId                 int64       `json:"userId"                 ` // 用户ID
	Type                   int         `json:"type"                   ` // 类型，0-个人，1-企业
	CompanyName            string      `json:"companyName"            ` // 公司名称
	BusinessNum            string      `json:"businessNum"            ` // 税号
	Name                   string      `json:"name"                   ` // 个人或法人姓名
	Idcard                 string      `json:"idcard"                 ` // 个人或法人身份证号
	Location               string      `json:"location"               ` // 省市区地址
	Address                string      `json:"address"                ` // 详细地址
	IdcardFrontPic         string      `json:"idcardFrontPic"         ` // 个人或法人身份证正面
	IdcardBackPic          string      `json:"idcardBackPic"          ` // 个人或法人身份证背面
	BusinessLicensePic     string      `json:"businessLicensePic"     ` // 营业执照图片
	Tag                    string      `json:"tag"                    ` // 标签，或业务类型，比如trtc,im
	GmtCreate              *gtime.Time `json:"gmtCreate"              ` // 创建时间
	GmtModify              *gtime.Time `json:"gmtModify"              ` // 修改时间
	BecomePractitionerTime *gtime.Time `json:"becomePractitionerTime" ` // 成为练习者时间
	IsDeleted              int         `json:"isDeleted"              ` // 是否删除，0-未删除，1-删除
	CompanyLogo            string      `json:"companyLogo"            ` // 账号头像
	Mobile                 string      `json:"mobile"                 ` // 登录手机号
	Mark                   string      `json:"mark"                   ` // 备注
	SettleCurrency         string      `json:"settleCurrency"         ` // 结算币种，null 代表结算用 CNY
	ServiceRate            int64       `json:"serviceRate"            ` // 服务费比例，万分位，百分比[0，10000)，精度为0.01%，如3即为0.03%
}
