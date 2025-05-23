package search

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type SearchReq struct {
	g.Meta    `path:"/key_search" tags:"Search" method:"get,post" summary:"Search"`
	SearchKey string `json:"searchKey" dc:"SearchKey, Will Search UserId|Email|UserName|CompanyName|SubscriptionId|VatNumber|InvoiceId||PaymentId" `
}

type PrecisionMatchObject struct {
	Type string      `json:"type" description:"match Type, user|invoice" `
	Id   string      `json:"id" description:"match Id user_id|invoice_id" `
	Data interface{} `json:"data" description:"match data" `
}

type SearchRes struct {
	PrecisionMatchObject *PrecisionMatchObject       `json:"precisionMatchObject" description:"PrecisionMatchObject" `
	UserAccounts         []*detail.UserAccountDetail `json:"matchUserAccounts" description:"MatchUserAccounts" `
	Invoices             []*bean.Invoice             `json:"matchInvoice" description:"MatchInvoice" `
}
