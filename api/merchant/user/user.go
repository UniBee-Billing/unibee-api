package user

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type ListReq struct {
	g.Meta    `path:"/list" tags:"User" method:"get,post" summary:"User List"`
	UserId    int    `json:"userId" dc:"Filter UserId" `
	FirstName string `json:"firstName" dc:"Search FirstName" `
	LastName  string `json:"lastName" dc:"Search LastName" `
	Email     string `json:"email" dc:"Search Filter Email" `
	Status    []int  `json:"status" dc:"Status, 0-Active｜2-Frozen" `
	//UserName           int    `json:"userName" dc:"Filter UserName, Default All" `
	//SubscriptionName   int    `json:"subscriptionName" dc:"Filter SubscriptionName, Default All" `
	//SubscriptionStatus int    `json:"subscriptionStatus" dc:"Filter SubscriptionStatus, Default All" `
	//PaymentMethod      int    `json:"paymentMethod" dc:"Filter GatewayDefaultPaymentMethod, Default All" `
	//BillingType        int    `json:"billingType" dc:"Filter BillingType, Default All" `
	DeleteInclude bool   `json:"deleteInclude" dc:"Deleted Involved，Need Admin" `
	SortField     string `json:"sortField" dc:"Sort，user_id|gmt_create|email|user_name|subscription_name|subscription_status|payment_method|recurring_amount|billing_type，Default gmt_create" `
	SortType      string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page          int    `json:"page"  dc:"Page,Start 0" `
	Count         int    `json:"count" dc:"Count OF Page" `
}

type ListRes struct {
	UserAccounts []*bean.UserAccountSimplify `json:"userAccounts" description:"UserAccounts" `
}

type GetReq struct {
	g.Meta `path:"/get" tags:"User" method:"get" summary:"Get User Profile"`
	UserId int64 `json:"userId" dc:"UserId" `
}

type GetRes struct {
	User *bean.UserAccountSimplify `json:"user" dc:"User"`
}

type FrozenReq struct {
	g.Meta `path:"/frozen_user" tags:"User" method:"post" summary:"Merchant Frozen User"`
	UserId int64 `json:"userId" dc:"UserId" `
}

type FrozenRes struct {
}

type ReleaseReq struct {
	g.Meta `path:"/release_user" tags:"User" method:"post" summary:"Merchant Release User"`
	UserId int64 `json:"userId" dc:"UserId" `
}

type ReleaseRes struct {
}

type SearchReq struct {
	g.Meta    `path:"/search" tags:"User" method:"get,post" summary:"User Search"`
	SearchKey string `json:"searchKey" dc:"SearchKey, Will Search UserId|Email|UserName|CompanyName|SubscriptionId|VatNumber|InvoiceId||PaymentId" `
}

type SearchRes struct {
	UserAccounts []*bean.UserAccountSimplify `json:"userAccounts" description:"UserAccounts" `
}

type UpdateReq struct {
	g.Meta          `path:"/update" tags:"User" method:"post" summary:"Update User Profile"`
	UserId          uint64 `json:"userId" dc:"User Id" v:"required"`
	FirstName       string `json:"firstName" dc:"First name" v:"required"`
	LastName        string `json:"lastName" dc:"Last Name" v:"required"`
	Email           string `json:"email" dc:"Email" v:"required"`
	Address         string `json:"address" dc:"Billing Address" v:"required"`
	CompanyName     string `json:"companyName" dc:"Company Name"`
	VATNumber       string `json:"vATNumber" dc:"VAT Number"`
	Phone           string `json:"phone" dc:"Phone"`
	Telegram        string `json:"telegram" dc:"Telegram"`
	WhatsApp        string `json:"whatsApp" dc:"WhatsApp"`
	WeChat          string `json:"weChat" dc:"WeChat"`
	LinkedIn        string `json:"LinkedIn" dc:"LinkedIn"`
	Facebook        string `json:"facebook" dc:"Facebook"`
	TikTok          string `json:"tiktok" dc:"Tiktok"`
	OtherSocialInfo string `json:"otherSocialInfo" dc:"Other Social Info"`
	PaymentMethod   string `json:"paymentMethod" dc:"Payment Method"`
	CountryCode     string `json:"countryCode" dc:"Country Code" v:"required"`
	CountryName     string `json:"countryName" dc:"Country Name" v:"required"`
}

type UpdateRes struct {
}
