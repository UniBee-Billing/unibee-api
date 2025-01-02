package session

import "github.com/gogf/gf/v2/frame/g"

type NewReq struct {
	g.Meta         `path:"/new_session" tags:"Session" method:"post" summary:"New Session" dc:"New session for user portal or web component"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId" v:"required"`
	Email          string `json:"email" dc:"Email" v:"required"`
	FirstName      string `json:"firstName" dc:"First Name"`
	LastName       string `json:"lastName" dc:"Last Name"`
	ReturnUrl      string `json:"returnUrl" dc:"ReturnUrl"`
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
