package user

import "github.com/gogf/gf/v2/frame/g"

type NewReq struct {
	g.Meta         `path:"/userPortal/new_session" tags:"User-Session-Controller" method:"post" summary:"New User Portal Session"`
	ExternalUserId string `p:"externalUserId" dc:"ExternalUserId" v:"required"`
	Email          string `p:"email" dc:"Email" v:"required"`
	FirstName      string `p:"firstName" dc:"First Name" v:"required"`
	LastName       string `p:"lastName" dc:"Last Name" v:"required"`
	ReturnUrl      string `p:"returnUrl" dc:"ReturnUrl"`
	Password       string `p:"password" dc:"Password"`
	Phone          string `p:"phone" dc:"Phone" `
	Address        string `p:"address" dc:"Address"`
}

type NewRes struct {
	UserId         string `json:"userId" dc:"UserId"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId"`
	Email          string `json:"email" dc:"Email"`
	Url            string `json:"url" dc:"Url"`
}
