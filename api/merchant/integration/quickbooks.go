package integration

import "github.com/gogf/gf/v2/frame/g"

type ConnectionQuickBooksReq struct {
	g.Meta    `path:"/quickbooks/get_authorization_url" tags:"Integrations" method:"get" summary:"Get Quickbooks Authorization Connection URL"`
	ReturnUrl string `json:"returnUrl" dc:"ReturnUrl"`
}

type ConnectionQuickBooksRes struct {
	AuthorizationURL string `json:"authorizationUrl" dc:"Authorization URL"`
}

type DisconnectionQuickBooksReq struct {
	g.Meta `path:"/quickbooks/disconnection" tags:"Integrations" method:"get" summary:"Disconnection Quickbooks"`
}

type DisconnectionQuickBooksRes struct {
}
