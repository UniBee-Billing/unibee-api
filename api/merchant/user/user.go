package user

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type NewReq struct {
	g.Meta         `path:"/new" tags:"User" method:"post" summary:"New User" dc:"User Creation If Not Exist "`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId"`
	Email          string `json:"email" dc:"Email" v:"required"`
	FirstName      string `json:"firstName" dc:"First Name"`
	LastName       string `json:"lastName" dc:"Last Name"`
	Password       string `json:"password" dc:"Password"`
	Phone          string `json:"phone" dc:"Phone" `
	Address        string `json:"address" dc:"Address"`
	Language       string `json:"language" dc:"Language"`
}

type NewRes struct {
	User *bean.UserAccount `json:"user" dc:"User Object"`
}

type ListReq struct {
	g.Meta          `path:"/list" tags:"User" method:"get,post" summary:"User List"`
	UserId          int64  `json:"userId" dc:"Filter UserId" `
	FirstName       string `json:"firstName" dc:"Search FirstName" `
	LastName        string `json:"lastName" dc:"Search LastName" `
	Email           string `json:"email" dc:"Search Filter Email" `
	PlanIds         []int  `json:"planIds" dc:"PlanIds, Search Filter PlanIds" `
	SubscriptionId  string `json:"subscriptionId" dc:"Search Filter SubscriptionId" `
	SubStatus       []int  `json:"subStatus" dc:"Filter, Default All，1-Pending｜2-Active｜3-Suspend | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete | 8-Processing | 9-Failed" `
	Status          []int  `json:"status" dc:"Status, 0-Active｜2-Frozen" `
	DeleteInclude   bool   `json:"deleteInclude" dc:"Deleted Involved，Need Admin" `
	SortField       string `json:"sortField" dc:"Sort，user_id|gmt_create|email|user_name|subscription_name|subscription_status|payment_method|recurring_amount|billing_type，Default gmt_create" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page,Start 0" `
	Count           int    `json:"count" dc:"Count OF Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
}

type ListRes struct {
	UserAccounts []*detail.UserAccountDetail `json:"userAccounts" description:"User Account Object List" `
	Total        int                         `json:"total" dc:"Total"`
}

type CountReq struct {
	g.Meta          `path:"/count" tags:"User" method:"get" summary:"User Count"`
	CreateTimeStart int64 `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64 `json:"createTimeEnd" dc:"CreateTimeEnd" `
}

type CountRes struct {
	Total int `json:"total" dc:"Total"`
}
type GetReq struct {
	g.Meta `path:"/get" tags:"User" method:"get" summary:"Get User Profile"`
	UserId int64 `json:"userId" dc:"UserId" `
}

type GetRes struct {
	User *detail.UserAccountDetail `json:"user" dc:"User"`
}

type FrozenReq struct {
	g.Meta `path:"/suspend_user" tags:"User" method:"post" summary:"Suspend User"`
	UserId int64 `json:"userId" dc:"UserId" `
}

type FrozenRes struct {
}

type ReleaseReq struct {
	g.Meta `path:"/resume_user" tags:"User" method:"post" summary:"Resume User"`
	UserId int64 `json:"userId" dc:"UserId" `
}

type ReleaseRes struct {
}

type SearchReq struct {
	g.Meta    `path:"/search" tags:"User" method:"get,post" summary:"User Search"`
	SearchKey string `json:"searchKey" dc:"SearchKey, Will Search UserId|Email|UserName|CompanyName|SubscriptionId|VatNumber|InvoiceId||PaymentId" `
}

type SearchRes struct {
	UserAccounts []*detail.UserAccountDetail `json:"userAccounts" description:"UserAccounts" `
}

type UpdateReq struct {
	g.Meta             `path:"/update" tags:"User" method:"post" summary:"Update User Profile"`
	UserId             *uint64                 `json:"userId" dc:"The id of user, either Email or UserId needed"`
	Email              *string                 `json:"email" dc:"The email of user, either Email or UserId needed"`
	FirstName          *string                 `json:"firstName" dc:"First name"`
	LastName           *string                 `json:"lastName" dc:"Last Name"`
	Address            *string                 `json:"address" dc:"Billing Address"`
	CompanyName        *string                 `json:"companyName" dc:"Company Name"`
	VATNumber          *string                 `json:"vATNumber" dc:"VAT Number"`
	RegistrationNumber *string                 `json:"registrationNumber" dc:"RegistrationNumber"`
	Phone              *string                 `json:"phone" dc:"Phone"`
	Telegram           *string                 `json:"telegram" dc:"Telegram"`
	WhatsApp           *string                 `json:"whatsApp" dc:"WhatsApp"`
	WeChat             *string                 `json:"weChat" dc:"WeChat"`
	LinkedIn           *string                 `json:"LinkedIn" dc:"LinkedIn"`
	Facebook           *string                 `json:"facebook" dc:"Facebook"`
	TikTok             *string                 `json:"tiktok" dc:"Tiktok"`
	OtherSocialInfo    *string                 `json:"otherSocialInfo" dc:"Other Social Info"`
	CountryCode        *string                 `json:"countryCode" dc:"Country Code"`
	CountryName        *string                 `json:"countryName" dc:"Country Name"`
	Type               *int64                  `json:"type" dc:"User type, 1-Individual|2-organization"`
	GatewayId          *uint64                 `json:"gatewayId" dc:"GatewayId"`
	GatewayPaymentType *string                 `json:"gatewayPaymentType" dc:"Gateway Payment Type"`
	PaymentMethodId    *string                 `json:"paymentMethodId" dc:"PaymentMethodId of gateway, available for card type gateway, payment automatic will enable if set" `
	City               *string                 `json:"city" dc:"city"`
	ZipCode            *string                 `json:"zipCode" dc:"zip_code"`
	Language           *string                 `json:"language" dc:"User Language, en|ru|cn|vi|bp"`
	ExternalUserId     *string                 `json:"externalUserId" dc:"ExternalUserId"`
	Metadata           *map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

type UpdateRes struct {
	User *detail.UserAccountDetail `json:"user" dc:"User"`
}

type ChangeGatewayReq struct {
	g.Meta             `path:"/change_gateway" tags:"User" method:"post" summary:"Change User Default Gateway" `
	UserId             uint64 `json:"userId" dc:"User Id" v:"required"`
	GatewayId          uint64 `json:"gatewayId" dc:"GatewayId" v:"required"`
	GatewayPaymentType string `json:"gatewayPaymentType" dc:"GatewayPaymentType"`
	GatewayUserId      string `json:"gatewayUserId" dc:"GatewayUserId, verify and save GatewayUserId via gateway"`
	PaymentMethodId    string `json:"paymentMethodId" dc:"PaymentMethodId of gateway, available for card type gateway, payment automatic will enable if set" `
}
type ChangeGatewayRes struct {
}

type ChangeEmailReq struct {
	g.Meta         `path:"/change_email" tags:"User" method:"post" summary:"Change User Email"`
	UserId         uint64 `json:"userId" dc:"The id of user, either ExternalUserId or UserId needed" `
	ExternalUserId string `json:"externalUserId" dc:"The externalUserId of user, either ExternalUserId or UserId needed"`
	NewEmail       string `json:"newEmail" dc:"Target Email want to change" v:"required"`
}

type ChangeEmailRes struct {
}

type ClearAutoChargeMethodReq struct {
	g.Meta `path:"/clear_auto_charge_method" tags:"User" method:"post" summary:"Clear AutoCharge Method"`
	UserId uint64 `json:"userId" dc:"The id of user" v:"required"`
}

type ClearAutoChargeMethodRes struct {
}

type NewAdminNoteReq struct {
	g.Meta `path:"/new_admin_note" tags:"User" method:"post" summary:"New Admin Note"`
	UserId uint64 `json:"userId" dc:"The id of user, either ExternalUserId or UserId needed" v:"required"`
	Note   string `json:"note" dc:"Note" v:"required"`
}

type NewAdminNoteRes struct {
}

type AdminNoteListReq struct {
	g.Meta `path:"/admin_note_list" tags:"User" method:"get,post" summary:"Get User Admin Note List"`
	UserId uint64 `json:"userId" dc:"The id of user, either ExternalUserId or UserId needed" v:"required"`
	Page   int    `json:"page"  dc:"Page, Start With 0" `
	Count  int    `json:"count" dc:"Count Of Page" `
}

type AdminNoteListRes struct {
	NoteLists []*detail.UserAdminNoteDetail `json:"noteLists"   description:""`
}
