package profile

import (
	entity "unibee/internal/model/entity/oversea_pay"

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

type PasswordResetReq struct {
	g.Meta      `path:"/passwordReset" tags:"Merchant-User-Profile-Controller" method:"post" summary:"Merchant User Reset Password"`
	OldPassword string `p:"oldPassword" dc:"OldPassword" v:"required"`
	NewPassword string `p:"newPassword" dc:"NewPassword" v:"required"`
}

type PasswordResetRes struct {
}
