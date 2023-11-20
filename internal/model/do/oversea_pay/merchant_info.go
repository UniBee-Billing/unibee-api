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
	g.Meta                 `orm:"table:merchant_info, do:true"`
	Id                     interface{} // 用户的ID
	CompanyId              interface{} // 公司ID
	UserId                 interface{} // 用户ID
	Type                   interface{} // 类型，0-个人，1-企业
	CompanyName            interface{} // 公司名称
	BusinessNum            interface{} // 税号
	Name                   interface{} // 个人或法人姓名
	Idcard                 interface{} // 个人或法人身份证号
	Location               interface{} // 省市区地址
	Address                interface{} // 详细地址
	IdcardFrontPic         interface{} // 个人或法人身份证正面
	IdcardBackPic          interface{} // 个人或法人身份证背面
	BusinessLicensePic     interface{} // 营业执照图片
	Tag                    interface{} // 标签，或业务类型，比如trtc,im
	GmtCreate              *gtime.Time // 创建时间
	GmtModify              *gtime.Time // 修改时间
	BecomePractitionerTime *gtime.Time // 成为练习者时间
	IsDeleted              interface{} // 是否删除，0-未删除，1-删除
	CompanyLogo            interface{} // 账号头像
	Mobile                 interface{} // 登录手机号
	Mark                   interface{} // 备注
	SettleCurrency         interface{} // 结算币种，null 代表结算用 CNY
	ServiceRate            interface{} // 服务费比例，万分位，百分比[0，10000)，精度为0.01%，如3即为0.03%
}
