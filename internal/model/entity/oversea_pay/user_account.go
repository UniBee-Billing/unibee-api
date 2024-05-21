// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAccount is the golang structure for table user_account.
type UserAccount struct {
	Id                 uint64      `json:"id"                 description:"userId"`                                                                                                    // userId
	GatewayId          string      `json:"gatewayId"          description:"gateway_id"`                                                                                                // gateway_id
	PaymentMethod      string      `json:"paymentMethod"      description:""`                                                                                                          //
	MerchantId         uint64      `json:"merchantId"         description:"merchant_id"`                                                                                               // merchant_id
	GmtCreate          *gtime.Time `json:"gmtCreate"          description:"create time"`                                                                                               // create time
	GmtModify          *gtime.Time `json:"gmtModify"          description:"update time"`                                                                                               // update time
	IsDeleted          int         `json:"isDeleted"          description:"0-UnDeleted，1-Deleted"`                                                                                     // 0-UnDeleted，1-Deleted
	Password           string      `json:"password"           description:"password , encrypt"`                                                                                        // password , encrypt
	UserName           string      `json:"userName"           description:"user name"`                                                                                                 // user name
	Mobile             string      `json:"mobile"             description:"mobile"`                                                                                                    // mobile
	Email              string      `json:"email"              description:"email"`                                                                                                     // email
	Gender             string      `json:"gender"             description:"gender"`                                                                                                    // gender
	AvatarUrl          string      `json:"avatarUrl"          description:"avator url"`                                                                                                // avator url
	ReMark             string      `json:"reMark"             description:"note"`                                                                                                      // note
	IsSpecial          int         `json:"isSpecial"          description:"is special account（0.no，1.yes）- deperated"`                                                                 // is special account（0.no，1.yes）- deperated
	Birthday           string      `json:"birthday"           description:"brithday"`                                                                                                  // brithday
	Profession         string      `json:"profession"         description:"profession"`                                                                                                // profession
	School             string      `json:"school"             description:"school"`                                                                                                    // school
	Custom             string      `json:"custom"             description:"custom"`                                                                                                    // custom
	LastLoginAt        int64       `json:"lastLoginAt"        description:"last login time, utc time"`                                                                                 // last login time, utc time
	IsRisk             int         `json:"isRisk"             description:"is risk account (deperated)"`                                                                               // is risk account (deperated)
	Version            int         `json:"version"            description:"version"`                                                                                                   // version
	Phone              string      `json:"phone"              description:"phone"`                                                                                                     // phone
	Address            string      `json:"address"            description:"address"`                                                                                                   // address
	FirstName          string      `json:"firstName"          description:"first name"`                                                                                                // first name
	LastName           string      `json:"lastName"           description:"last name"`                                                                                                 // last name
	CompanyName        string      `json:"companyName"        description:"company name"`                                                                                              // company name
	VATNumber          string      `json:"vATNumber"          description:"vat number"`                                                                                                // vat number
	Telegram           string      `json:"telegram"           description:"telegram"`                                                                                                  // telegram
	WhatsAPP           string      `json:"whatsAPP"           description:"whats app"`                                                                                                 // whats app
	WeChat             string      `json:"weChat"             description:"wechat"`                                                                                                    // wechat
	TikTok             string      `json:"tikTok"             description:"tictok"`                                                                                                    // tictok
	LinkedIn           string      `json:"linkedIn"           description:"linkedin"`                                                                                                  // linkedin
	Facebook           string      `json:"facebook"           description:"facebook"`                                                                                                  // facebook
	OtherSocialInfo    string      `json:"otherSocialInfo"    description:""`                                                                                                          //
	CountryCode        string      `json:"countryCode"        description:"country_code"`                                                                                              // country_code
	CountryName        string      `json:"countryName"        description:"country_name"`                                                                                              // country_name
	SubscriptionName   string      `json:"subscriptionName"   description:"subscription name"`                                                                                         // subscription name
	SubscriptionId     string      `json:"subscriptionId"     description:"subscription id"`                                                                                           // subscription id
	SubscriptionStatus int         `json:"subscriptionStatus" description:"sub status，0-Init | 1-Pending｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete"` // sub status，0-Init | 1-Pending｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete
	RecurringAmount    int64       `json:"recurringAmount"    description:"total recurring amount, cent"`                                                                              // total recurring amount, cent
	BillingType        int         `json:"billingType"        description:"1-recurring,2-one-time"`                                                                                    // 1-recurring,2-one-time
	TimeZone           string      `json:"timeZone"           description:""`                                                                                                          //
	CreateTime         int64       `json:"createTime"         description:"create utc time"`                                                                                           // create utc time
	ExternalUserId     string      `json:"externalUserId"     description:"external_user_id"`                                                                                          // external_user_id
	Status             int         `json:"status"             description:"0-Active, 2-Suspend"`                                                                                       // 0-Active, 2-Suspend
}
