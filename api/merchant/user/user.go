package user

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

type ListReq struct {
	g.Meta     `path:"/user_list" tags:"Merchant-User-Controller" method:"post" summary:"User List"`
	MerchantId uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
	UserId     int    `p:"userId" dc:"Filter UserId" `
	FirstName  string `p:"firstName" dc:"Search FirstName" `
	LastName   string `p:"lastName" dc:"Search LastName" `
	Email      string `p:"email" dc:"Search Filter Email" `
	Status     []int  `p:"status" dc:"Status, 0-Active｜2-Frozen" `
	//UserName           int    `p:"userName" dc:"Filter UserName, Default All" `
	//SubscriptionName   int    `p:"subscriptionName" dc:"Filter SubscriptionName, Default All" `
	//SubscriptionStatus int    `p:"subscriptionStatus" dc:"Filter SubscriptionStatus, Default All" `
	//PaymentMethod      int    `p:"paymentMethod" dc:"Filter GatewayDefaultPaymentMethod, Default All" `
	//BillingType        int    `p:"billingType" dc:"Filter BillingType, Default All" `
	DeleteInclude bool   `p:"deleteInclude" dc:"Deleted Involved，Need Admin" `
	SortField     string `p:"sortField" dc:"Sort，user_id|gmt_create|email|user_name|subscription_name|subscription_status|payment_method|recurring_amount|billing_type，Default gmt_create" `
	SortType      string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page          int    `p:"page"  dc:"Page,Start 0" `
	Count         int    `p:"count" dc:"Count OF Page" `
}

type ListRes struct {
	UserAccounts []*entity.UserAccount `json:"userAccounts" description:"UserAccounts" `
}

type GetReq struct {
	g.Meta `path:"/get_user_profile" tags:"Merchant-User-Controller" method:"get" summary:"Get User Profile"`
	UserId int64 `p:"userId" dc:"UserId" `
}

type GetRes struct {
	User *entity.UserAccount `p:"user" dc:"User"`
}

type FrozenReq struct {
	g.Meta `path:"/frozen_user" tags:"Merchant-User-Controller" method:"get" summary:"Merchant Frozen User"`
	UserId int64 `p:"userId" dc:"UserId" `
}

type FrozenRes struct {
}

type ReleaseReq struct {
	g.Meta `path:"/release_user" tags:"Merchant-User-Controller" method:"get" summary:"Merchant Release User"`
	UserId int64 `p:"userId" dc:"UserId" `
}

type ReleaseRes struct {
}

type SearchReq struct {
	g.Meta     `path:"/user_search" tags:"Merchant-User-Controller" method:"post" summary:"User Search"`
	MerchantId uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
	SearchKey  string `p:"searchKey" dc:"SearchKey, Will Search UserId|Email|UserName|CompanyName|SubscriptionId|VatNumber|InvoiceId||PaymentId" `
}

type SearchRes struct {
	UserAccounts []*entity.UserAccount `json:"userAccounts" description:"UserAccounts" `
}

type UserProfileUpdateReq struct {
	g.Meta          `path:"/update_user_profile" tags:"Merchant-User-Controller" method:"post" summary:"Update User Profile"`
	UserId          uint64 `p:"userId" dc:"User Id" v:"required"`
	FirstName       string `p:"firstName" dc:"First name" v:"required"`
	LastName        string `p:"lastName" dc:"Last Name" v:"required"`
	Email           string `p:"email" dc:"Email" v:"required"`
	Address         string `p:"address" dc:"Billing Address" v:"required"`
	CompanyName     string `p:"companyName" dc:"Company Name"`
	VATNumber       string `p:"vATNumber" dc:"VAT Number"`
	Phone           string `p:"phone" dc:"Phone"`
	Telegram        string `p:"telegram" dc:"Telegram"`
	WhatsApp        string `p:"whatsApp" dc:"WhatsApp"`
	WeChat          string `p:"weChat" dc:"WeChat"`
	LinkedIn        string `p:"LinkedIn" dc:"LinkedIn"`
	Facebook        string `p:"facebook" dc:"Facebook"`
	TikTok          string `p:"tiktok" dc:"Tiktok"`
	OtherSocialInfo string `p:"otherSocialInfo" dc:"Other Social Info"`
	PaymentMethod   string `p:"paymentMethod" dc:"Payment Method"`
	CountryCode     string `p:"countryCode" dc:"Country Code" v:"required"`
	CountryName     string `p:"countryName" dc:"Country Name" v:"required"`
}

// with token to be implemented in the future
type UserProfileUpdateRes struct {
	//User *entity.UserAccount `p:"user" dc:"User"`
	// Token string `p:"token" dc:"token string"`
}
