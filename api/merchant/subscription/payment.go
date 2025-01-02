package subscription

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/api/merchant/payment"
)

type NewPaymentReq struct {
	g.Meta            `path:"/payment/new" tags:"Subscription Payment" method:"post" summary:"New Subscription Payment"`
	ExternalPaymentId string                 `json:"externalPaymentId" dc:"ExternalPaymentId should unique for payment"`
	ExternalUserId    string                 `json:"externalUserId" dc:"ExternalUserId, unique, either ExternalUserId&Email or UserId needed"`
	Email             string                 `json:"email" dc:"Email, either ExternalUserId&Email or UserId needed"`
	UserId            uint64                 `json:"userId" dc:"UserId, either ExternalUserId&Email or UserId needed"`
	Currency          string                 `json:"currency" dc:"Currency, either Currency&TotalAmount or PlanId needed" `
	TotalAmount       int64                  `json:"totalAmount" dc:"Total PaymentAmount, Cent, either TotalAmount&Currency or PlanId needed"`
	PlanId            uint64                 `json:"planId" dc:"PlanId, either TotalAmount&Currency or PlanId needed"`
	GatewayId         uint64                 `json:"gatewayId"   dc:"GatewayId" v:"required"`
	RedirectUrl       string                 `json:"redirectUrl" dc:"Redirect Url"`
	CancelUrl         string                 `json:"cancelUrl" dc:"CancelUrl"`
	CountryCode       string                 `json:"countryCode" dc:"CountryCode"`
	Name              string                 `json:"name" dc:"Name"`
	Description       string                 `json:"description" dc:"Description"`
	Items             []*payment.Item        `json:"items" dc:"Items"`
	Metadata          map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
	GasPayer          string                 `json:"gasPayer" dc:"who pay the gas, merchant|user"`
}

type NewPaymentRes struct {
	Status            int         `json:"status" dc:"Status, 10-Created|20-Success|30-Failed|40-Cancelled"`
	PaymentId         string      `json:"paymentId" dc:"The unique id of payment"`
	ExternalPaymentId string      `json:"externalPaymentId" dc:"The external unique id of payment"`
	Link              string      `json:"link"`
	Action            *gjson.Json `json:"action" dc:"action"`
}

type OnetimeAddonNewReq struct {
	g.Meta             `path:"/new_onetime_addon_payment" tags:"Subscription Payment" method:"post" summary:"New Subscription Onetime Addon Payment" dc:"Create payment for subscription onetime addon purchase"`
	SubscriptionId     string                 `json:"subscriptionId" dc:"SubscriptionId, id of subscription which addon will attached, either SubscriptionId or UserId needed, The only one active subscription of userId will attach the addon"`
	UserId             uint64                 `json:"userId" dc:"UserId, either SubscriptionId or UserId needed, The only one active subscription will update if userId provide instead of subscriptionId"`
	AddonId            uint64                 `json:"addonId" dc:"AddonId, id of one-time addon, the new payment will created base on the addon's amount'" v:"required"`
	Quantity           int64                  `json:"quantity" dc:"Quantity, quantity of the new payment which one-time addon purchased"  v:"required"`
	ReturnUrl          string                 `json:"returnUrl"  dc:"ReturnUrl, the addon's payment will redirect based on the returnUrl provided when it's back from gateway side"  `
	Metadata           map[string]interface{} `json:"metadata" dc:"Metadata，custom data"`
	DiscountCode       string                 `json:"discountCode" dc:"DiscountCode"`
	DiscountAmount     *int64                 `json:"discountAmount"     dc:"Amount of discount"`
	DiscountPercentage *int64                 `json:"discountPercentage" dc:"Percentage of discount, 100=1%, ignore if discountAmount provide"`
	GatewayId          *uint64                `json:"gatewayId" dc:"GatewayId, use user's gateway if not provide"`
	TaxPercentage      *int64                 `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, use subscription's taxPercentage if not provide"`
}

type OnetimeAddonNewRes struct {
	SubscriptionOnetimeAddon *bean.SubscriptionOnetimeAddon `json:"subscriptionOnetimeAddon" dc:"SubscriptionOnetimeAddon, object of onetime-addon purchased"`
	Paid                     bool                           `json:"paid" dc:"true|false,automatic payment is default behavior for one-time addon purchased, payment will create attach to the purchase, when payment is success, return false, otherwise false"`
	Link                     string                         `json:"link" dc:"if automatic payment is false, Gateway Link will provided that manual payment needed"`
	Invoice                  *bean.Invoice                  `json:"invoice" dc:"invoice of one-time payment"`
}

type OnetimeAddonListReq struct {
	g.Meta `path:"/onetime_addon_list" tags:"Subscription Payment" method:"get" summary:"Get Subscription Onetime Addon List"`
	UserId uint64 `json:"userId" dc:"UserId" v:"required"`
	Page   int    `json:"page"  dc:"Page, Start With 0" `
	Count  int    `json:"count" dc:"Count Of Page" `
}

type OnetimeAddonListRes struct {
	SubscriptionOnetimeAddons []*detail.SubscriptionOnetimeAddonDetail `json:"subscriptionOnetimeAddons" description:"SubscriptionOnetimeAddons" `
}
