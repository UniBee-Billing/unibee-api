package payment

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type NewReq struct {
	g.Meta            `path:"/new" tags:"Payment" method:"post" summary:"New Payment"`
	ExternalPaymentId string            `json:"externalPaymentId" dc:"ExternalPaymentId should unique for payment" v:"required"`
	ExternalUserId    string            `json:"externalUserId" dc:"ExternalUserId, should unique for user" v:"required"`
	Email             string            `json:"email" dc:"Email" v:"required"`
	Currency          string            `json:"currency" dc:"Currency"  v:"required"`
	TotalAmount       int64             `json:"totalAmount" dc:"Total PaymentAmount, Cent" v:"required"`
	GatewayId         uint64            `json:"gatewayId"   dc:"GatewayId" v:"required"`
	RedirectUrl       string            `json:"redirectUrl" dc:"Redirect Url"`
	CountryCode       string            `json:"countryCode" dc:"CountryCode"`
	Items             []*Item           `json:"items" dc:"Items"`
	Metadata          map[string]string `json:"metadata" dc:"Metadata，Map"`
}

type NewRes struct {
	Status            int         `json:"status" dc:"Status, 10-Created|20-Success|30-Failed|40-Cancelled"`
	PaymentId         string      `json:"paymentId" dc:"PaymentId"`
	ExternalPaymentId string      `json:"externalPaymentId" dc:"ExternalPaymentId"`
	Link              string      `json:"link"`
	Action            *gjson.Json `json:"action" dc:"action"`
}

type Item struct {
	Currency               string `json:"currency"`
	Amount                 int64  `json:"amount" dc:"the item total amount,cent"`
	Tax                    int64  `json:"tax"`
	AmountExcludingTax     int64  `json:"amountExcludingTax"`
	TaxScale               int64  `json:"taxScale" dc:"Tax Scale，1000 = 10%"`
	UnitAmountExcludingTax int64  `json:"unitAmountExcludingTax"`
	Description            string `json:"description"`
	Quantity               int64  `json:"quantity"`
}

type DetailReq struct {
	g.Meta    `path:"/detail" tags:"Payment" method:"get" summary:"Query Payment Detail"`
	PaymentId string `json:"paymentId" dc:"PaymentId" v:"required"`
}
type DetailRes struct {
	PaymentDetail *ro.PaymentDetailRo `json:"paymentDetail" dc:"PaymentDetail"`
}

type ListReq struct {
	g.Meta      `path:"/list" tags:"Payment" method:"get" summary:"Query Payment List"`
	GatewayId   uint64 `json:"gatewayId"   dc:"GatewayId"`
	UserId      int64  `json:"userId" dc:"UserId " `
	Email       string `json:"email" dc:"Email"`
	Status      int    `json:"status" dc:"Status, 10-Created|20-Success|30-Failed|40-Cancelled"`
	Currency    string `json:"currency" dc:"Currency"`
	CountryCode string `json:"countryCode" dc:"CountryCode"`
	SortField   string `json:"sortField" dc:"Sort Field，user_id|create_time|status" `
	SortType    string `json:"sortType" dc:"Sort Type，asc|desc" `
	Page        int    `json:"page"  dc:"Page, Start With 0" `
	Count       int    `json:"count"  dc:"Count" dc:"Count Of Page" `
}

type ListRes struct {
	PaymentDetails []*ro.PaymentDetailRo `json:"paymentDetails" dc:"PaymentDetails"`
}
