package detail

import (
	"context"
	"strconv"
	"unibee/api/bean"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

type UserAccountDetail struct {
	Id                 uint64                `json:"id"                 description:"userId"`                                                                                                           // userId
	MerchantId         uint64                `json:"merchantId"         description:"merchant_id"`                                                                                                      // merchant_id
	UserName           string                `json:"userName"           description:"user name"`                                                                                                        // user name
	Mobile             string                `json:"mobile"             description:"mobile"`                                                                                                           // mobile
	Email              string                `json:"email"              description:"email"`                                                                                                            // email
	Gender             string                `json:"gender"             description:"gender"`                                                                                                           // gender
	AvatarUrl          string                `json:"avatarUrl"          description:"avator url"`                                                                                                       // avator url
	ReMark             string                `json:"reMark"             description:"note"`                                                                                                             // note
	IsSpecial          int                   `json:"isSpecial"          description:"is special account（0.no，1.yes）- deperated"`                                                                        // is special account（0.no，1.yes）- deperated
	Birthday           string                `json:"birthday"           description:"brithday"`                                                                                                         // brithday
	Profession         string                `json:"profession"         description:"profession"`                                                                                                       // profession
	School             string                `json:"school"             description:"school"`                                                                                                           // school
	Custom             string                `json:"custom"             description:"custom"`                                                                                                           // custom
	LastLoginAt        int64                 `json:"lastLoginAt"        description:"last login time, utc time"`                                                                                        // last login time, utc time
	IsRisk             int                   `json:"isRisk"             description:"is risk account (deperated)"`                                                                                      // is risk account (deperated)
	GatewayId          uint64                `json:"gatewayId"          description:"gateway_id"`                                                                                                       // gateway_id
	Version            int                   `json:"version"            description:"version"`                                                                                                          // version
	Phone              string                `json:"phone"              description:"phone"`                                                                                                            // phone
	Address            string                `json:"address"            description:"address"`                                                                                                          // address
	FirstName          string                `json:"firstName"          description:"first name"`                                                                                                       // first name
	LastName           string                `json:"lastName"           description:"last name"`                                                                                                        // last name
	CompanyName        string                `json:"companyName"        description:"company name"`                                                                                                     // company name
	VATNumber          string                `json:"vATNumber"          description:"vat number"`                                                                                                       // vat number
	Telegram           string                `json:"telegram"           description:"telegram"`                                                                                                         // telegram
	WhatsAPP           string                `json:"whatsAPP"           description:"whats app"`                                                                                                        // whats app
	WeChat             string                `json:"weChat"             description:"wechat"`                                                                                                           // wechat
	TikTok             string                `json:"tikTok"             description:"tictok"`                                                                                                           // tictok
	LinkedIn           string                `json:"linkedIn"           description:"linkedin"`                                                                                                         // linkedin
	Facebook           string                `json:"facebook"           description:"facebook"`                                                                                                         // facebook
	OtherSocialInfo    string                `json:"otherSocialInfo"    description:""`                                                                                                                 //
	PaymentMethod      string                `json:"paymentMethod"      description:""`                                                                                                                 //
	CountryCode        string                `json:"countryCode"        description:"country_code"`                                                                                                     // country_code
	CountryName        string                `json:"countryName"        description:"country_name"`                                                                                                     // country_name
	SubscriptionName   string                `json:"subscriptionName"   description:"subscription name"`                                                                                                // subscription name
	SubscriptionId     string                `json:"subscriptionId"     description:"subscription id"`                                                                                                  // subscription id
	SubscriptionStatus int                   `json:"subscriptionStatus" description:"sub status， 1-Pending｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete | 8-Processing"` // sub status，0-Init | 1-Create｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete
	RecurringAmount    int64                 `json:"recurringAmount"    description:"total recurring amount, cent"`                                                                                     // total recurring amount, cent
	BillingType        int                   `json:"billingType"        description:"1-recurring,2-one-time"`                                                                                           // 1-recurring,2-one-time
	TimeZone           string                `json:"timeZone"           description:""`                                                                                                                 //
	CreateTime         int64                 `json:"createTime"         description:"create utc time"`                                                                                                  // create utc time
	ExternalUserId     string                `json:"externalUserId"     description:"external_user_id"`                                                                                                 // external_user_id
	Status             int                   `json:"status"             description:"0-Active, 2-Suspend"`
	TaxPercentage      int64                 `json:"taxPercentage"      description:"taxPercentage，1000 = 10%"`               // taxPercentage，1000 = 10%
	Type               int64                 `json:"type"               description:"User type, 1-Individual|2-organization"` // User type, 1-Individual|2-organization
	Gateway            *bean.GatewaySimplify `json:"gateway"            description:"Gateway"`
	City               string                `json:"city" dc:"city"`
	ZipCode            string                `json:"zipCode" dc:"zip_code"`
}

func ConvertUserAccountToDetail(ctx context.Context, one *entity.UserAccount) *UserAccountDetail {
	if one == nil {
		return nil
	}
	gatewayId, _ := strconv.ParseUint(one.GatewayId, 10, 64)
	return &UserAccountDetail{
		Id:                 one.Id,
		MerchantId:         one.MerchantId,
		UserName:           one.UserName,
		Mobile:             one.Mobile,
		Email:              one.Email,
		Gender:             one.Gender,
		Type:               one.Type,
		TaxPercentage:      one.TaxPercentage,
		AvatarUrl:          one.AvatarUrl,
		ReMark:             one.ReMark,
		IsSpecial:          one.IsSpecial,
		Birthday:           one.Birthday,
		Profession:         one.Profession,
		School:             one.School,
		Custom:             one.Custom,
		LastLoginAt:        one.LastLoginAt,
		IsRisk:             one.IsRisk,
		GatewayId:          gatewayId,
		Version:            one.Version,
		Phone:              one.Phone,
		Address:            one.Address,
		FirstName:          one.FirstName,
		LastName:           one.LastName,
		CompanyName:        one.CompanyName,
		VATNumber:          one.VATNumber,
		Telegram:           one.Telegram,
		WhatsAPP:           one.WhatsAPP,
		WeChat:             one.WeChat,
		TikTok:             one.TikTok,
		LinkedIn:           one.LinkedIn,
		Facebook:           one.Facebook,
		OtherSocialInfo:    one.OtherSocialInfo,
		PaymentMethod:      one.PaymentMethod,
		CountryCode:        one.CountryCode,
		CountryName:        one.CountryName,
		SubscriptionName:   one.SubscriptionName,
		SubscriptionId:     one.SubscriptionId,
		SubscriptionStatus: one.SubscriptionStatus,
		RecurringAmount:    one.RecurringAmount,
		BillingType:        one.BillingType,
		TimeZone:           one.TimeZone,
		CreateTime:         one.CreateTime,
		ExternalUserId:     one.ExternalUserId,
		Status:             one.Status,
		City:               one.City,
		ZipCode:            one.ZipCode,
		Gateway:            bean.SimplifyGateway(query.GetGatewayById(ctx, gatewayId)),
	}
}
