package profile

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type GetReq struct {
	g.Meta `path:"/get" tags:"User-Profile" method:"get" summary:"Get User Profile"`
}

type GetRes struct {
	User *bean.UserAccountSimplify `json:"user" dc:"User"`
}

type LogoutReq struct {
	g.Meta `path:"/logout" tags:"User-Profile" method:"post" summary:"User Logout"`
}

type LogoutRes struct {
}

type UpdateReq struct {
	g.Meta          `path:"/update" tags:"User-Profile" method:"post" summary:"Update User Profile"`
	FirstName       string `json:"firstName" dc:"First name"`
	LastName        string `json:"lastName" dc:"Last Name"`
	Email           string `json:"email" dc:"Email" v:"required"`
	Address         string `json:"address" dc:"Billing Address" v:"required"`
	CompanyName     string `json:"companyName" dc:"Company Name"`
	VATNumber       string `json:"vATNumber" dc:"VAT Number"`
	Phone           string `json:"phone" dc:"Phone"`
	Telegram        string `json:"telegram" dc:"Telegram"`
	WhatsApp        string `json:"WhatsApp" dc:"WhatsApp"`
	WeChat          string `json:"WeChat" dc:"WeChat"`
	LinkedIn        string `json:"LinkedIn" dc:"LinkedIn"`
	Facebook        string `json:"facebook" dc:"Facebook"`
	TikTok          string `json:"tiktok" dc:"Tiktok"`
	OtherSocialInfo string `json:"otherSocialInfo" dc:"Other Social Info"`
	TimeZone        string `json:"timeZone" dc:"User TimeZone"`
	CountryCode     string `json:"countryCode" dc:"Country Code" v:"required"`
	CountryName     string `json:"countryName" dc:"Country Name" v:"required"`
}

type UpdateRes struct {
	User *bean.UserAccountSimplify `json:"user" dc:"User"`
}

type PasswordResetReq struct {
	g.Meta      `path:"/passwordReset" tags:"User-Profile" method:"post" summary:"User Reset Password"`
	OldPassword string `json:"oldPassword" dc:"OldPassword" v:"required"`
	NewPassword string `json:"newPassword" dc:"NewPassword" v:"required"`
}

type PasswordResetRes struct {
}

type ChangeGatewayReq struct {
	g.Meta          `path:"/change_gateway" tags:"User" method:"post" summary:"ChangeUserDefaultGateway" `
	GatewayId       uint64 `json:"gatewayId" dc:"GatewayId" v:"required"`
	PaymentMethodId string `json:"paymentMethodId" dc:"PaymentMethodId of gateway, available for card type gateway, payment automatic will enable if set" `
}
type ChangeGatewayRes struct {
}
