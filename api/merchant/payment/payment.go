package payment

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type NewReq struct {
	g.Meta            `path:"/new" tags:"Payment" method:"post" summary:"New Payment"`
	ExternalPaymentId string            `json:"externalPaymentId" dc:"ExternalPaymentId should unique for payment" v:"required"`
	ExternalUserId    string            `json:"externalUserId" dc:"ExternalUserId, should unique for user" v:"required"`
	Email             string            `json:"email" dc:"Email" v:"required"`
	Currency          string            `json:"currency" dc:"Currency"  v:"required"`
	TotalAmount       int64             `json:"totalAmount" dc:"Total PaymentAmount, Cent" v:"required"`
	RedirectUrl       string            `json:"redirectUrl" dc:"Redirect Url" v:"required"`
	CountryCode       string            `json:"countryCode" dc:"CountryCode" v:"required"`
	GatewayId         uint64            `json:"gatewayId"   dc:"GatewayId" v:"required"`
	Items             []*Item           `json:"lineItems" dc:"Items"`
	Metadata          map[string]string `json:"reference" dc:"Metadata，Map" v:""`
}
type NewRes struct {
	Status            int         `json:"status" dc:"Status, 10-Created|20-Success|30-Failed|40-Cancelled"`
	PaymentId         string      `json:"paymentId" dc:"PaymentId"`
	ExternalPaymentId string      `json:"externalPaymentId" dc:"ExternalPaymentId"`
	Action            *gjson.Json `json:"action" dc:"action"`
}

type Item struct {
	UnitAmountExcludingTax int64  `json:"unitAmountExcludingTax" dc:"UnitAmountExcludingTax" v:"required"`
	Quantity               int64  `json:"quantity" dc:"Quantity" v:"required"`
	Description            string `json:"description" dc:"Description" v:""`
	TaxScale               int64  `json:"taxScale" dc:"TaxScale" v:"required"`
	ProductUrl             string `json:"productUrl" dc:"ProductUrl"`
	ImageUrl               string `json:"imageUrl" dc:"ImageUrl"`
}

type DetailReq struct {
	g.Meta    `path:"/detail" tags:"Payment" method:"get" summary:"Query Payment Detail"`
	PaymentId string `json:"paymentId" dc:"PaymentId" v:"required"`
}
type DetailRes struct {
}

type ListReq struct {
	g.Meta         `path:"/list" tags:"Payment" method:"get" summary:"Query Payment List"`
	GatewayId      uint64 `json:"gatewayId"   dc:"GatewayId"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId"`
	Email          string `json:"email" dc:"Email"`
	Currency       string `json:"currency" dc:"Currency"`
	CountryCode    string `json:"countryCode" dc:"CountryCode"`
}
type ListRes struct {
}
