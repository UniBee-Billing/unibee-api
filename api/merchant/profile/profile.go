package profile

import (
	entity "go-oversea-pay/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/frame/g"
)

type ProfileReq struct {
	g.Meta `path:"/profile" tags:"Merchant-User-Profile-Controller" method:"get" summary:"Get Merchant User Profile"`
	// Email string `p:"email" dc:"email" v:"required"`
	// Password  string `p:"password" dc:"password" v:"required"`
}

type ProfileRes struct {
	MerchantUser *entity.MerchantUserAccount `p:"merchantUser" dc:"Merchant User"`
	// Token string `p:"token" dc:"token string"`
}

type LogoutReq struct {
	g.Meta `path:"/user_logout" tags:"Merchant-User-Profile-Controller" method:"post" summary:"User Logout"`
}

type LogoutRes struct {
}
