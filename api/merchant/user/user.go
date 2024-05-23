package user

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type NewReq struct {
	g.Meta         `path:"/new" tags:"User" method:"post" summary:"NewUser" dc:"New User"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId"`
	Email          string `json:"email" dc:"Email" v:"required"`
	FirstName      string `json:"firstName" dc:"First Name"`
	LastName       string `json:"lastName" dc:"Last Name"`
	Password       string `json:"password" dc:"Password"`
	Phone          string `json:"phone" dc:"Phone" `
	Address        string `json:"address" dc:"Address"`
}

type NewRes struct {
	User *bean.UserAccountSimplify `json:"user" dc:"User Object"`
}

type ListReq struct {
	g.Meta        `path:"/list" tags:"User" method:"get,post" summary:"UserList"`
	UserId        int    `json:"userId" dc:"Filter UserId" `
	FirstName     string `json:"firstName" dc:"Search FirstName" `
	LastName      string `json:"lastName" dc:"Search LastName" `
	Email         string `json:"email" dc:"Search Filter Email" `
	Status        []int  `json:"status" dc:"Status, 0-Active｜2-Frozen" `
	DeleteInclude bool   `json:"deleteInclude" dc:"Deleted Involved，Need Admin" `
	SortField     string `json:"sortField" dc:"Sort，user_id|gmt_create|email|user_name|subscription_name|subscription_status|payment_method|recurring_amount|billing_type，Default gmt_create" `
	SortType      string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page          int    `json:"page"  dc:"Page,Start 0" `
	Count         int    `json:"count" dc:"Count OF Page" `
}

type ListRes struct {
	UserAccounts []*detail.UserAccountDetail `json:"userAccounts" description:"User Account Object List" `
	Total        int                         `json:"total" dc:"Total"`
}

type GetReq struct {
	g.Meta `path:"/get" tags:"User" method:"get" summary:"GetUserProfile"`
	UserId int64 `json:"userId" dc:"UserId" `
}

type GetRes struct {
	User *detail.UserAccountDetail `json:"user" dc:"User"`
}

type FrozenReq struct {
	g.Meta `path:"/suspend_user" tags:"User" method:"post" summary:"SuspendUser"`
	UserId int64 `json:"userId" dc:"UserId" `
}

type FrozenRes struct {
}

type ReleaseReq struct {
	g.Meta `path:"/resume_user" tags:"User" method:"post" summary:"ResumeUser"`
	UserId int64 `json:"userId" dc:"UserId" `
}

type ReleaseRes struct {
}

type SearchReq struct {
	g.Meta    `path:"/search" tags:"User" method:"get,post" summary:"UserSearch"`
	SearchKey string `json:"searchKey" dc:"SearchKey, Will Search UserId|Email|UserName|CompanyName|SubscriptionId|VatNumber|InvoiceId||PaymentId" `
}

type SearchRes struct {
	UserAccounts []*detail.UserAccountDetail `json:"userAccounts" description:"UserAccounts" `
}

type UpdateReq struct {
	g.Meta          `path:"/update" tags:"User" method:"post" summary:"UpdateUserProfile"`
	UserId          uint64  `json:"userId" dc:"User Id" v:"required"`
	FirstName       string  `json:"firstName" dc:"First name" v:"required"`
	LastName        string  `json:"lastName" dc:"Last Name" v:"required"`
	Email           string  `json:"email" dc:"Email" v:"required"`
	Address         string  `json:"address" dc:"Billing Address" v:"required"`
	CompanyName     string  `json:"companyName" dc:"Company Name"`
	VATNumber       *string `json:"vATNumber" dc:"VAT Number"`
	Phone           string  `json:"phone" dc:"Phone"`
	Telegram        string  `json:"telegram" dc:"Telegram"`
	WhatsApp        string  `json:"whatsApp" dc:"WhatsApp"`
	WeChat          string  `json:"weChat" dc:"WeChat"`
	LinkedIn        string  `json:"LinkedIn" dc:"LinkedIn"`
	Facebook        string  `json:"facebook" dc:"Facebook"`
	TikTok          string  `json:"tiktok" dc:"Tiktok"`
	OtherSocialInfo string  `json:"otherSocialInfo" dc:"Other Social Info"`
	CountryCode     *string `json:"countryCode" dc:"Country Code"`
	CountryName     *string `json:"countryName" dc:"Country Name"`
	Type            *int64  `json:"type" dc:"User type, 1-Individual|2-organization"`
	GatewayId       *uint64 `json:"gatewayId" dc:"GatewayId"`
	PaymentMethodId *string `json:"paymentMethodId" dc:"PaymentMethodId of gateway, available for card type gateway, payment automatic will enable if set" `
}

type UpdateRes struct {
}

type ChangeGatewayReq struct {
	g.Meta          `path:"/change_gateway" tags:"User" method:"post" summary:"ChangeUserDefaultGateway" `
	UserId          uint64 `json:"userId" dc:"User Id" v:"required"`
	GatewayId       uint64 `json:"gatewayId" dc:"GatewayId" v:"required"`
	PaymentMethodId string `json:"paymentMethodId" dc:"PaymentMethodId of gateway, available for card type gateway, payment automatic will enable if set" `
}
type ChangeGatewayRes struct {
}
