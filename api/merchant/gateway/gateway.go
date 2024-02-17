package gateway

import "github.com/gogf/gf/v2/frame/g"

type CheckAndSetupReq struct {
	g.Meta `path:"/gateway_check_and_setup" tags:"Merchant-Gateway-Controller" method:"post" summary:"Gateway Check And Setup"`
}
type CheckAndSetupRes struct {
}
