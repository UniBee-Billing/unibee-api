package session

import "github.com/gogf/gf/v2/frame/g"

type NewReq struct {
	g.Meta         `path:"/new_session" tags:"Session" method:"post" summary:"New User Portal Session" dc:"New session for user portal or web component"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId" v:"required"`
	Email          string `json:"email" dc:"Email" v:"required"`
	FirstName      string `json:"firstName" dc:"First Name"`
	LastName       string `json:"lastName" dc:"Last Name"`
	ReturnUrl      string `json:"returnUrl" dc:"ReturnUrl"`
	CancelUrl      string `json:"cancelUrl" dc:"CancelUrl"`
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
	ReturnUrl      string `json:"returnUrl"  dc:"ReturnUrl"`
	CancelUrl      string `json:"cancelUrl" dc:"CancelUrl"`
}

type NewSubUpdatePageRes struct {
	Url string `json:"url" dc:"Url"`
}
