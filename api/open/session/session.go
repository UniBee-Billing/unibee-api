package session

import "github.com/gogf/gf/v2/frame/g"

type UserPortalSessionReq struct {
	g.Meta         `path:"/userPortal/session" tags:"Open-Session-Controller" method:"post" summary:"Create User Portal Session"`
	ExternalUserId string `p:"externalUserId" dc:"ExternalUserId" v:"required"`
	Email          int64  `p:"email" dc:"Email" v:"required"`
	ReturnUrl      string `p:"returnUrl" dc:"ReturnUrl"`
}

type UserPortalSessionRes struct {
	UserId string `json:"userId" dc:"UserId"`
	Url    string `json:"url" dc:"Url"`
}
