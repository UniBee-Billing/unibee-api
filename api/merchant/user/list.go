package user

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type ListReq struct {
	g.Meta             `path:"/user_list" tags:"Merchant-User-Controller" method:"post" summary:"User List"`
	MerchantId         int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	UserId             int    `p:"userId" dc:"Filter UserId, Default All" `
	Email              int    `p:"email" dc:"Filter Email, Default All" `
	UserName           int    `p:"userName" dc:"Filter UserName, Default All" `
	SubscriptionName   int    `p:"subscriptionName" dc:"Filter SubscriptionName, Default All" `
	SubscriptionStatus int    `p:"subscriptionStatus" dc:"Filter SubscriptionStatus, Default All" `
	PaymentMethod      int    `p:"paymentMethod" dc:"Filter PaymentMethod, Default All" `
	BillingType        int    `p:"billingType" dc:"Filter BillingType, Default All" `
	DeleteInclude      bool   `p:"deleteInclude" dc:"Deleted Involved，Need Admin" `
	SortField          string `p:"sortField" dc:"Sort，user_id|gmt_create|email|user_name|subscription_name|subscription_status|payment_method|recurring_amount|billing_type，Default gmt_create" `
	SortType           string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page               int    `p:"page"  dc:"Page,Start 0" `
	Count              int    `p:"count" dc:"Count OF Page" `
}

type ListRes struct {
	UserAccounts []*entity.UserAccount `json:"userAccounts" description:"UserAccounts" `
}

type SearchReq struct {
	g.Meta     `path:"/user_search" tags:"Merchant-User-Controller" method:"post" summary:"User Search"`
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	SearchKey  string `p:"searchKey" dc:"SearchKey, Will Search UserId|Email|UserName|CompanyName|VatNumber|InvoiceId|SubscriptionId|PaymentId" `
}

type SearchRes struct {
	UserAccounts []*entity.UserAccount `json:"userAccounts" description:"UserAccounts" `
}
