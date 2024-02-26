package search

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "unibee/internal/model/entity/oversea_pay"
)

type SearchReq struct {
	g.Meta    `path:"/key_search" tags:"Merchant-Search-Controller" method:"post" summary:"Merchant Search"`
	SearchKey string `p:"searchKey" dc:"SearchKey, Will Search UserId|Email|UserName|CompanyName|SubscriptionId|VatNumber|InvoiceId||PaymentId" `
}

type PrecisionMatchObject struct {
	Type string      `json:"type" description:"match Type, user|invoice" `
	Id   string      `json:"id" description:"match Id user_id|invoice_id" `
	Data interface{} `json:"data" description:"match data" `
}

type SearchRes struct {
	PrecisionMatchObject *PrecisionMatchObject `json:"precisionMatchObject" description:"PrecisionMatchObject" `
	UserAccounts         []*entity.UserAccount `json:"matchUserAccounts" description:"MatchUserAccounts" `
	Invoices             []*entity.Invoice     `json:"matchInvoice" description:"MatchInvoice" `
}
