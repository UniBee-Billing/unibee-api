package profile

import (
	entity "unibee-api/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/frame/g"
)

type ProfileReq struct {
	g.Meta `path:"/profile" tags:"User-Profile-Controller" method:"get" summary:"Get User Profile"`
	// Email string `p:"email" dc:"email" v:"required"`
	// Password  string `p:"password" dc:"password" v:"required"`
}

// with token to be implemented in the future
type ProfileRes struct {
	User *entity.UserAccount `p:"user" dc:"User"`
	// Token string `p:"token" dc:"token string"`
}

type LogoutReq struct {
	g.Meta `path:"/user_logout" tags:"User-Profile-Controller" method:"post" summary:"User Logout"`
}

type LogoutRes struct {
}

type ProfileUpdateReq struct {
	g.Meta          `path:"/profile" tags:"User-Profile-Controller" method:"post" summary:"Update User Profile"`
	Id              uint64 `p:"id" dc:"User Id" v:"required"`
	FirstName       string `p:"firstName" dc:"First name" v:"required"`
	LastName        string `p:"lastName" dc:"Last Name" v:"required"`
	Email           string `p:"email" dc:"Email" v:"required"`
	Address         string `p:"address" dc:"Billing Address" v:"required"`
	CompanyName     string `p:"companyName" dc:"Company Name"`
	VATNumber       string `p:"vATNumber" dc:"VAT Number"`
	Phone           string `p:"phone" dc:"Phone"`
	Telegram        string `p:"telegram" dc:"Telegram"`
	WhatsApp        string `p:"WhatsApp" dc:"WhatsApp"`
	WeChat          string `p:"WeChat" dc:"WeChat"`
	LinkedIn        string `p:"LinkedIn" dc:"LinkedIn"`
	Facebook        string `p:"facebook" dc:"Facebook"`
	TikTok          string `p:"tiktok" dc:"Tiktok"`
	OtherSocialInfo string `p:"otherSocialInfo" dc:"Other Social Info"`
	PaymentMethod   string `p:"paymentMethod" dc:"Payment Method"`
	TimeZone        string `p:"timeZone" dc:"User TimeZone"`
	CountryCode     string `p:"countryCode" dc:"Country Code" v:"required"`
	CountryName     string `p:"countryName" dc:"Country Name" v:"required"`
}

type ProfileUpdateRes struct {
	User *entity.UserAccount `p:"user" dc:"User"`
}

type PasswordResetReq struct {
	g.Meta      `path:"/passwordReset" tags:"User-Profile-Controller" method:"post" summary:"User Reset Password"`
	OldPassword string `p:"oldPassword" dc:"OldPassword" v:"required"`
	NewPassword string `p:"newPassword" dc:"NewPassword" v:"required"`
}

type PasswordResetRes struct {
}
