package payment

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type NewReq struct {
	g.Meta            `path:"/new" tags:"Payment" method:"post" summary:"NewPayment"`
	ExternalPaymentId string                 `json:"externalPaymentId" dc:"ExternalPaymentId should unique for payment"`
	ExternalUserId    string                 `json:"externalUserId" dc:"ExternalUserId, unique, either ExternalUserId&Email or UserId needed"`
	Email             string                 `json:"email" dc:"Email, either ExternalUserId&Email or UserId needed"`
	UserId            uint64                 `json:"userId" dc:"UserId, either ExternalUserId&Email or UserId needed"`
	Currency          string                 `json:"currency" dc:"Currency, either Currency&TotalAmount or PlanId needed" `
	TotalAmount       int64                  `json:"totalAmount" dc:"Total PaymentAmount, Cent, either TotalAmount&Currency or PlanId needed"`
	PlanId            uint64                 `json:"planId" dc:"PlanId, either TotalAmount&Currency or PlanId needed"`
	GatewayId         uint64                 `json:"gatewayId"   dc:"GatewayId" v:"required"`
	RedirectUrl       string                 `json:"redirectUrl" dc:"Redirect Url"`
	CountryCode       string                 `json:"countryCode" dc:"CountryCode"`
	Description       string                 `json:"description" dc:"Description"`
	Items             []*Item                `json:"items" dc:"Items"`
	Metadata          map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
	GasPayer          string                 `json:"gasPayer" dc:"who pay the gas, merchant|user"`
}

type NewRes struct {
	Status            int         `json:"status" dc:"Status, 10-Created|20-Success|30-Failed|40-Cancelled"`
	PaymentId         string      `json:"paymentId" dc:"The unique id of payment"`
	ExternalPaymentId string      `json:"externalPaymentId" dc:"The external unique id of payment"`
	Link              string      `json:"link"`
	Action            *gjson.Json `json:"action" dc:"action"`
}

type Item struct {
	Amount                 int64  `json:"amount" dc:"item total amount, sum(item.amount) should equal to totalAmount, cent"  v:"required"`
	Description            string `json:"description" dc:"The item description of payment" v:"required" `
	Quantity               int64  `json:"quantity" dc:"The item quantity of payment" `
	UnitAmountExcludingTax int64  `json:"unitAmountExcludingTax"`
	Currency               string `json:"currency" dc:"The item currency of payment" `
	Tax                    int64  `json:"tax" dc:"tax = amount - amountExcludingTax"`
	AmountExcludingTax     int64  `json:"amountExcludingTax" dc:"amountExcludingTax = unitAmountExcludingTax * quantity"`
	TaxPercentage          int64  `json:"taxPercentage" dc:"The tax percentage of payment，1000 = 10%"`
}

type DetailReq struct {
	g.Meta    `path:"/detail" tags:"Payment" method:"get" summary:"PaymentDetail"`
	PaymentId string `json:"paymentId" dc:"The unique id of payment" v:"required"`
}
type DetailRes struct {
	PaymentDetail *detail.PaymentDetail `json:"paymentDetail" dc:"Payment Detail Object"`
}

type ListReq struct {
	g.Meta      `path:"/list" tags:"Payment" method:"get" summary:"PaymentList"`
	GatewayId   uint64 `json:"gatewayId"   dc:"The filter unique id of gateway"`
	UserId      uint64 `json:"userId" dc:"The filter userid of payment" `
	Email       string `json:"email" dc:"The filter email of payment"`
	Status      int    `json:"status" dc:"The filter status of payment, 10-Created|20-Success|30-Failed|40-Cancelled"`
	Currency    string `json:"currency" dc:"The filter currency of payment"`
	CountryCode string `json:"countryCode" dc:"The filter country code of payment"`
	SortField   string `json:"sortField" dc:"Sort Field，user_id|create_time|status" `
	SortType    string `json:"sortType" dc:"Sort Type，asc|desc" `
	Page        int    `json:"page"  dc:"Page, Start With 0" `
	Count       int    `json:"count"  dc:"Count" dc:"Count Of Page" `
}

type ListRes struct {
	PaymentDetails []*detail.PaymentDetail `json:"paymentDetails" dc:"Payment Detail Object List"`
	Total          int                     `json:"total" dc:"Total"`
}
