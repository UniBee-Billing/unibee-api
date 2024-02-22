package payment

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type NewPaymentReq struct {
	g.Meta                   `path:"/new_payment" tags:"OneTime-Payment-Controller" method:"post" summary:"New Payment"`
	MerchantId               uint64          `p:"merchantId" dc:"MerchantId" v:"required"`
	MerchantPaymentId        string          `p:"merchantPaymentId" dc:"MerchantPaymentId" v:"required"`
	TotalAmount              *AmountVo       `json:"totalAmount" dc:"Total Amount, Cent" v:"required"`
	PaymentMethod            *MethodListReq  `json:"paymentMethod"   in:"query" dc:"Payment Method" v:"required"`
	RedirectUrl              string          `p:"redirectUrl" dc:"Redirect Url" v:"required"`
	CountryCode              string          `p:"countryCode" dc:"CountryCode" v:"required"`
	ShopperLocale            string          `p:"shopperLocale" dc:"Shopper Locale" v:"required"`
	ShopperEmail             string          `p:"shopperEmail" dc:"Shopper Email" v:"required"`
	ShopperUserId            string          `p:"shopperUserId" dc:"shopper Id, Unique" v:"required"`
	LineItems                []*OutLineItem  `p:"lineItems" dc:"LineItems" v:"required"`
	DeviceType               string          `p:"deviceType" dc:"DeviceType,Android|iOS|Web"`
	Platform                 string          `p:"platform" dc:"Platform（WEB，WAP，APP, MINI, WALLET）"`
	ShopperIP                string          `p:"shopperIP" dc:"Shopper IP（v4，v6）"`
	TelephoneNumber          string          `p:"telephoneNumber" dc:"TelephoneNumber"`
	BrowserInfo              string          `p:"browserInfo" dc:"browserInfo" v:""`
	ShopperInteraction       string          `p:"shopperInteraction" dc:"ShopperInteraction" v:""`
	RecurringProcessingToken string          `p:"recurringProcessingToken" dc:"RecurringProcessingToken" v:""`
	ShopperName              *OutShopperName `p:"shopperName" dc:"shopperName" v:""`
	//BillingAddress           *OutPayAddress     `p:"billingAddress" dc:"账单地址" v:""`
	//DetailAddress            *OutPayAddress     `p:"detailAddress" dc:"邮寄地址" v:""`
	Capture                bool              `p:"capture" dc:"Capture Immediate" v:""`
	CaptureDelayHours      int               `p:"captureDelayHours" dc:"Delay Capture Hours" v:""`
	MerchantOrderReference string            `p:"merchantOrderReference" dc:"Deprecated" v:""`
	Metadata               map[string]string `p:"reference" dc:"Metadata，Map" v:""`
	DateOfBrith            string            `p:"dateOfBrith" dc:"DateOfBrith，Format YYYY-MM-DD" v:""`
}
type NewPaymentRes struct {
	Status            string      `p:"status" dc:"Status"`
	PaymentId         string      `p:"paymentId" dc:"PaymentId"`
	MerchantPaymentId string      `p:"merchantPaymentId" dc:"MerchantPaymentId"`
	Action            *gjson.Json `p:"action" dc:"action"`
}

type OutShopperName struct {
	FirstName string `p:"firstName" dc:"First Name" v:"required"`
	LastName  string `p:"lastName" dc:"Last Name" v:"required"`
	Gender    string `p:"gender" dc:"Gender" v:"required"`
}

type OutPayAddress struct {
	City              string `p:"city" dc:"City" v:"required"`
	Country           string `p:"country" dc:"Country" v:"required"`
	HouseNumberOrName string `p:"houseNumberOrName" dc:"HouseNumberOrName" v:"required"`
	PostalCode        string `p:"postalCode" dc:"PostalCode" v:"required"`
	StateOrProvince   string `p:"stateOrProvince" dc:"StateOrProvince" v:"required"`
	Street            string `p:"street" dc:"Street" v:"required"`
}

type OutLineItem struct {
	UnitAmountExcludingTax int64  `p:"unitAmountExcludingTax" dc:"UnitAmountExcludingTax" v:"required"`
	Quantity               int64  `p:"quantity" dc:"Quantity" v:"required"`
	Description            string `p:"description" dc:"Description" v:""`
	TaxScale               int64  `p:"taxScale" dc:"TaxScale" v:"required"`
	ProductUrl             string `p:"productUrl" dc:"ProductUrl"`
	ImageUrl               string `p:"imageUrl" dc:"ImageUrl"`
	// todo mark discount need
}

type MethodListReq struct {
	g.Meta  `path:"/paymentMethodList" tags:"OneTime-Payment-Controller" method:"post" summary:"Payment Method Query (Support Klarna、Evonet）"`
	TokenId string `p:"tokenId" dc:"TokenId" v:""`
	Gateway string `p:"type" dc:"Gateway" v:"required"`
}
type MethodListRes struct {
}

type MethodIssur struct {
	Name     string `p:"name" dc:"Name" v:""`
	Id       string `p:"id" dc:"Method Id" v:""`
	Disabled string `p:"disabled" dc:"" v:""`
}

type DetailReq struct {
	g.Meta    `path:"/paymentDetail/{PaymentId}" tags:"OneTime-Payment-Controller" method:"post" summary:"Query Payment Detail"`
	PaymentId string `in:"path" dc:"PaymentId" v:"required"`
}
type DetailRes struct {
}
