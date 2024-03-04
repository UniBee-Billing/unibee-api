package payment

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type NewPaymentReq struct {
	g.Meta                   `path:"/new_payment" tags:"OneTime-Payment" method:"post" summary:"New Payment"`
	MerchantPaymentId        string          `json:"merchantPaymentId" dc:"MerchantPaymentId" v:"required"`
	TotalAmount              *AmountVo       `json:"totalAmount" dc:"Total Amount, Cent" v:"required"`
	PaymentMethod            *MethodListReq  `json:"paymentMethod"   in:"query" dc:"Payment Method" v:"required"`
	RedirectUrl              string          `json:"redirectUrl" dc:"Redirect Url" v:"required"`
	CountryCode              string          `json:"countryCode" dc:"CountryCode" v:"required"`
	ShopperLocale            string          `json:"shopperLocale" dc:"Shopper Locale" v:"required"`
	ShopperEmail             string          `json:"shopperEmail" dc:"Shopper Email" v:"required"`
	ShopperUserId            string          `json:"shopperUserId" dc:"shopper Id, Unique" v:"required"`
	LineItems                []*OutLineItem  `json:"lineItems" dc:"LineItems" v:"required"`
	DeviceType               string          `json:"deviceType" dc:"DeviceType,Android|iOS|Web"`
	Platform                 string          `json:"platform" dc:"Platform（WEB，WAP，APP, MINI, WALLET）"`
	ShopperIP                string          `json:"shopperIP" dc:"Shopper IP（v4，v6）"`
	TelephoneNumber          string          `json:"telephoneNumber" dc:"TelephoneNumber"`
	BrowserInfo              string          `json:"browserInfo" dc:"browserInfo" v:""`
	ShopperInteraction       string          `json:"shopperInteraction" dc:"ShopperInteraction" v:""`
	RecurringProcessingToken string          `json:"recurringProcessingToken" dc:"RecurringProcessingToken" v:""`
	ShopperName              *OutShopperName `json:"shopperName" dc:"shopperName" v:""`
	//BillingAddress           *OutPayAddress     `json:"billingAddress" dc:"账单地址" v:""`
	//DetailAddress            *OutPayAddress     `json:"detailAddress" dc:"邮寄地址" v:""`
	Capture                bool              `json:"capture" dc:"Capture Immediate" v:""`
	CaptureDelayHours      int               `json:"captureDelayHours" dc:"Delay Capture Hours" v:""`
	MerchantOrderReference string            `json:"merchantOrderReference" dc:"Deprecated" v:""`
	Metadata               map[string]string `json:"reference" dc:"Metadata，Map" v:""`
	DateOfBrith            string            `json:"dateOfBrith" dc:"DateOfBrith，Format YYYY-MM-DD" v:""`
}
type NewPaymentRes struct {
	Status            string      `json:"status" dc:"Status"`
	PaymentId         string      `json:"paymentId" dc:"PaymentId"`
	MerchantPaymentId string      `json:"merchantPaymentId" dc:"MerchantPaymentId"`
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
	// todo mark discount need
}

type MethodListReq struct {
	g.Meta  `path:"/paymentMethodList" tags:"OneTime-Payment" method:"post" summary:"Payment Method Query (Support Klarna、Evonet）"`
	TokenId string `json:"tokenId" dc:"TokenId" v:""`
	Gateway string `json:"type" dc:"Gateway" v:"required"`
}
type MethodListRes struct {
}

type MethodIssur struct {
	Name     string `json:"name" dc:"Name" v:""`
	Id       string `json:"id" dc:"Method Id" v:""`
	Disabled string `json:"disabled" dc:"" v:""`
}

type DetailReq struct {
	g.Meta    `path:"/paymentDetail/{PaymentId}" tags:"OneTime-Payment" method:"post" summary:"Query Payment Detail"`
	PaymentId string `in:"path" dc:"PaymentId" v:"required"`
}
type DetailRes struct {
}
