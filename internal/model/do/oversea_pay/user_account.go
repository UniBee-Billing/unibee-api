// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAccount is the golang structure of table user_account for DAO operations like Where/Data.
type UserAccount struct {
	g.Meta             `orm:"table:user_account, do:true"`
	Id                 interface{} // userId
	Email              interface{} // email
	GatewayId          interface{} // gateway_id
	PaymentMethod      interface{} //
	CountryCode        interface{} // country_code
	CountryName        interface{} // country_name
	VATNumber          interface{} // vat number
	TaxPercentage      interface{} // taxPercentage，1000 = 10%
	Type               interface{} // User type, 1-Individual|2-organization
	MerchantId         interface{} // merchant_id
	GmtCreate          *gtime.Time // create time
	GmtModify          *gtime.Time // update time
	IsDeleted          interface{} // 0-UnDeleted，1-Deleted
	Password           interface{} // password , encrypt
	UserName           interface{} // user name
	Mobile             interface{} // mobile
	Gender             interface{} // gender
	AvatarUrl          interface{} // avator url
	ReMark             interface{} // note
	IsSpecial          interface{} // is special account（0.no，1.yes）- deperated
	Birthday           interface{} // brithday
	Profession         interface{} // profession
	School             interface{} // school
	Custom             interface{} // custom
	LastLoginAt        interface{} // last login time, utc time
	IsRisk             interface{} // is risk account (deperated)
	Version            interface{} // version
	Phone              interface{} // phone
	Address            interface{} // address
	FirstName          interface{} // first name
	LastName           interface{} // last name
	CompanyName        interface{} // company name
	Telegram           interface{} // telegram
	WhatsAPP           interface{} // whats app
	WeChat             interface{} // wechat
	TikTok             interface{} // tictok
	LinkedIn           interface{} // linkedin
	Facebook           interface{} // facebook
	OtherSocialInfo    interface{} //
	SubscriptionName   interface{} // subscription name
	PlanId             interface{} // PlanId
	SubscriptionId     interface{} // subscription id
	SubscriptionStatus interface{} // sub status，0-Init | 1-Pending｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete
	RecurringAmount    interface{} // total recurring amount, cent
	BillingType        interface{} // 1-recurring,2-one-time
	TimeZone           interface{} //
	CreateTime         interface{} // create utc time
	ExternalUserId     interface{} // external_user_id
	Status             interface{} // 0-Active, 2-Suspend
	City               interface{} // city
	ZipCode            interface{} // zip_code
}
