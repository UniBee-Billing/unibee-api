package profile

import (
	entity "go-oversea-pay/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/frame/g"
)

type ProfileReq struct {
	g.Meta `path:"/profile" tags:"User-Profile-Controller" method:"get" summary:"get user profile"`
	// Email string `p:"email" dc:"email" v:"required"`
	// Password  string `p:"password" dc:"password" v:"required"`
}

// with token to be implemented in the future
type ProfileRes struct {
	User *entity.UserAccount `p:"user" dc:"user"`
	// Token string `p:"token" dc:"token string"`
}

type ProfileUpdateReq struct {
	g.Meta          `path:"/profile" tags:"User-Profile-Controller" method:"post" summary:"update user profile"`
	Id              uint64 `p:"id" dc:"user id" v:"required"`
	FirstName       string `p:"firstName" dc:"first name" v:"required"`
	LastName        string `p:"lastName" dc:"last name" v:"required"`
	Email           string `p:"email" dc:"email" v:"required"`
	Address         string `p:"address" dc:"billing address" v:"required"`
	CompanyName     string `p:"companyName" dc:"company name"`
	VATNumber       string `p:"vATNumber" dc:"VAT number"`
	Phone           string `p:"phone" dc:"phone"`
	Telegram        string `p:"telegram" dc:"telegram"`
	WhatsApp        string `p:"WhatsApp" dc:"whatsApp"`
	WeChat          string `p:"WeChat" dc:"weChat"`
	LinkedIn        string `p:"LinkedIn" dc:"linkedIn"`
	Facebook        string `p:"facebook" dc:"facebook"`
	TikTok          string `p:"tiktok" dc:"tiktok"`
	OtherSocialInfo string `p:"otherSocialInfo" dc:"other social info"`
	PaymmentMethod  string `p:"paymentMethod" dc:"payment method"`
	// Email string `p:"email" dc:"email" v:"required"`
	// Password  string `p:"password" dc:"password" v:"required"`
}

// with token to be implemented in the future
type ProfileUpdateRes struct {
	User *entity.UserAccount `p:"user" dc:"user"`
	// Token string `p:"token" dc:"token string"`
}
