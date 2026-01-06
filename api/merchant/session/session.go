package session

import "github.com/gogf/gf/v2/frame/g"

type NewReq struct {
	g.Meta         `path:"/new_session" tags:"Session" method:"post" summary:"New Checkout Session" dc:"New session for hosted checkout or client portal. You can create user and get ClientSession from here, then append it to the checkout link (copied from Admin Portal Plan) as a query parameter, e.g. https://cs.unibee.dev/hosted/checkout?planId=253&env=prod&session=${clientSession}"`
	Email          string `json:"email" dc:"Email" v:"required"`
	ReturnUrl      string `json:"returnUrl" dc:"ReturnUrl, back to returnUrl if checkout completed"`
	CancelUrl      string `json:"cancelUrl" dc:"CancelUrl, back to cancelUrl if checkout cancelled"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId"`
	FirstName      string `json:"firstName" dc:"First Name"`
	LastName       string `json:"lastName" dc:"Last Name"`
	Password       string `json:"password" dc:"Password"`
	Phone          string `json:"phone" dc:"Phone" `
	Address        string `json:"address" dc:"Address"`
}

type NewRes struct {
	UserId         string `json:"userId" dc:"UserId"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId"`
	Email          string `json:"email" dc:"Email"`
	Url            string `json:"url" dc:"Url"`
	ClientToken    string `json:"clientToken" dc:"ClientToken"`
	ClientSession  string `json:"clientSession" dc:"ClientSession"`
}

type NewSubUpdatePageReq struct {
	g.Meta         `path:"/user_sub_update_url" tags:"Session" method:"get,post" summary:"Get User Subscription Update Page Url"`
	Email          string `json:"email" dc:"Email" dc:"Email, unique, either ExternalUserId&Email or UserId needed"`
	UserId         uint64 `json:"userId" dc:"UserId" dc:"UserId, unique, either ExternalUserId&Email or UserId needed"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId, unique, either ExternalUserId&Email or UserId needed"`
	ProductId      int64  `json:"productId" dc:"Id of product" dc:"default product will use if productId not specified"`
	PlanId         int64  `json:"planId" dc:"Id of plan to update" dc:"Id of plan to update"`
	VatCountryCode string `json:"vatCountryCode" dc:"Vat Country Code"`
	ReturnUrl      string `json:"returnUrl"  dc:"ReturnUrl"`
	CancelUrl      string `json:"cancelUrl" dc:"CancelUrl"`
}

type NewSubUpdatePageRes struct {
	Url string `json:"url" dc:"Url"`
}
