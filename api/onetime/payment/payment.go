package payment

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type NewPaymentReq struct {
	g.Meta            `path:"/new_payment" tags:"OneTime-Payment" method:"post" summary:"New Payment"`
	ExternalPaymentId string            `json:"externalPaymentId" dc:"ExternalPaymentId should unique for payment" v:"required"`
	ExternalUserId    string            `json:"externalUserId" dc:"ExternalUserId, should unique for user" v:"required"`
	Email             string            `json:"email" dc:"Email" v:"required"`
	TotalAmount       *AmountVo         `json:"totalAmount" dc:"Total Amount, Cent" v:"required"`
	RedirectUrl       string            `json:"redirectUrl" dc:"Redirect Url" v:"required"`
	CountryCode       string            `json:"countryCode" dc:"CountryCode" v:"required"`
	PaymentMethod     *MethodListReq    `json:"paymentMethod"   in:"query" dc:"Payment Method" v:"required"`
	LineItems         []*OutLineItem    `json:"lineItems" dc:"LineItems" v:"required"`
	Metadata          map[string]string `json:"reference" dc:"Metadata，Map" v:""`
}
type NewPaymentRes struct {
	Status            string      `json:"status" dc:"Status"`
	PaymentId         string      `json:"paymentId" dc:"PaymentId"`
	ExternalPaymentId string      `json:"externalPaymentId" dc:"ExternalPaymentId"`
	Action            *gjson.Json `json:"action" dc:"action"`
}

type OutShopperName struct {
	FirstName string `json:"firstName" dc:"First Name" v:"required"`
	LastName  string `json:"lastName" dc:"Last Name" v:"required"`
	Gender    string `json:"gender" dc:"Gender" v:"required"`
}

type OutPayAddress struct {
	City              string `json:"city" dc:"City" v:"required"`
	Country           string `json:"country" dc:"Country" v:"required"`
	HouseNumberOrName string `json:"houseNumberOrName" dc:"HouseNumberOrName" v:"required"`
	PostalCode        string `json:"postalCode" dc:"PostalCode" v:"required"`
	StateOrProvince   string `json:"stateOrProvince" dc:"StateOrProvince" v:"required"`
	Street            string `json:"street" dc:"Street" v:"required"`
}

type OutLineItem struct {
	UnitAmountExcludingTax int64  `json:"unitAmountExcludingTax" dc:"UnitAmountExcludingTax" v:"required"`
	Quantity               int64  `json:"quantity" dc:"Quantity" v:"required"`
	Description            string `json:"description" dc:"Description" v:""`
	TaxScale               int64  `json:"taxScale" dc:"TaxScale" v:"required"`
	ProductUrl             string `json:"productUrl" dc:"ProductUrl"`
	ImageUrl               string `json:"imageUrl" dc:"ImageUrl"`
}

type MethodListReq struct {
	g.Meta  `path:"/paymentMethodList" tags:"OneTime-Payment" method:"post" summary:"Payment Method Query (Support Klarna、Evonet）"`
	Gateway string `json:"type" dc:"Gateway" v:"required"`
}
type MethodListRes struct {
}

type DetailReq struct {
	g.Meta    `path:"/paymentDetail/{PaymentId}" tags:"OneTime-Payment" method:"post" summary:"Query Payment Detail"`
	PaymentId string `in:"path" dc:"PaymentId" v:"required"`
}
type DetailRes struct {
}
