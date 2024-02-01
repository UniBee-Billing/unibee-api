// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAccount is the golang structure for table user_account.
type UserAccount struct {
	Id                 uint64      `json:"id"                 description:"userId"`                                                                                                   // userId
	GmtCreate          *gtime.Time `json:"gmtCreate"          description:"create time"`                                                                                              // create time
	GmtModify          *gtime.Time `json:"gmtModify"          description:"修改时间"`                                                                                                     // 修改时间
	IsDeleted          int         `json:"isDeleted"          description:"0-UnDeleted，1-Deleted"`                                                                                    // 0-UnDeleted，1-Deleted
	Password           string      `json:"password"           description:"密码，加密存储"`                                                                                                  // 密码，加密存储
	UserName           string      `json:"userName"           description:"用户名"`                                                                                                      // 用户名
	Mobile             string      `json:"mobile"             description:"手机号"`                                                                                                      // 手机号
	Email              string      `json:"email"              description:"邮箱"`                                                                                                       // 邮箱
	Gender             string      `json:"gender"             description:"性别"`                                                                                                       // 性别
	AvatarUrl          string      `json:"avatarUrl"          description:"头像url"`                                                                                                    // 头像url
	ReMark             string      `json:"reMark"             description:"备注"`                                                                                                       // 备注
	IsSpecial          int         `json:"isSpecial"          description:"是否是特殊账号（0.否，1.是）"`                                                                                         // 是否是特殊账号（0.否，1.是）
	Birthday           string      `json:"birthday"           description:"生日"`                                                                                                       // 生日
	Profession         string      `json:"profession"         description:"职业"`                                                                                                       // 职业
	School             string      `json:"school"             description:"学校"`                                                                                                       // 学校
	Custom             string      `json:"custom"             description:"其他"`                                                                                                       // 其他
	NearTime           *gtime.Time `json:"nearTime"           description:"最近登录时间"`                                                                                                   // 最近登录时间
	Wid                string      `json:"wid"                description:"盟有wid"`                                                                                                    // 盟有wid
	IsRisk             int         `json:"isRisk"             description:"风控：0.低风险，1.中风险，2.高风险"`                                                                                     // 风控：0.低风险，1.中风险，2.高风险
	Channel            string      `json:"channel"            description:"渠道"`                                                                                                       // 渠道
	Version            int         `json:"version"            description:"版本"`                                                                                                       // 版本
	Phone              string      `json:"phone"              description:""`                                                                                                         //
	Address            string      `json:"address"            description:""`                                                                                                         //
	FirstName          string      `json:"firstName"          description:""`                                                                                                         //
	LastName           string      `json:"lastName"           description:""`                                                                                                         //
	CompanyName        string      `json:"companyName"        description:""`                                                                                                         //
	VATNumber          string      `json:"vATNumber"          description:""`                                                                                                         //
	Telegram           string      `json:"telegram"           description:""`                                                                                                         //
	WhatsAPP           string      `json:"whatsAPP"           description:""`                                                                                                         //
	WeChat             string      `json:"weChat"             description:""`                                                                                                         //
	TikTok             string      `json:"tikTok"             description:""`                                                                                                         //
	LinkedIn           string      `json:"linkedIn"           description:""`                                                                                                         //
	Facebook           string      `json:"facebook"           description:""`                                                                                                         //
	OtherSocialInfo    string      `json:"otherSocialInfo"    description:""`                                                                                                         //
	PaymentMethod      string      `json:"paymentMethod"      description:""`                                                                                                         //
	CountryCode        string      `json:"countryCode"        description:"country_code"`                                                                                             // country_code
	CountryName        string      `json:"countryName"        description:"country_name"`                                                                                             // country_name
	SubscriptionName   string      `json:"subscriptionName"   description:""`                                                                                                         //
	SubscriptionId     int64       `json:"subscriptionId"     description:""`                                                                                                         //
	SubscriptionStatus int         `json:"subscriptionStatus" description:"sub status，0-Init | 1-Create｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete"` // sub status，0-Init | 1-Create｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete
	RecurringAmount    int64       `json:"recurringAmount"    description:""`                                                                                                         //
	BillingType        int         `json:"billingType"        description:"1-recurring,2-one-time"`                                                                                   // 1-recurring,2-one-time
}
