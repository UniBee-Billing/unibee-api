package member

import (
	entity "unibee/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/frame/g"
)

type ProfileReq struct {
	g.Meta `path:"/profile" tags:"Merchant-Member-Profile-Controller" method:"get" summary:"Get Merchant Member Profile"`
}

type ProfileRes struct {
	MerchantMember *entity.MerchantMember `json:"merchantMember" dc:"Merchant Member"`
}

type LogoutReq struct {
	g.Meta `path:"/user_logout" tags:"Merchant-Member-Profile-Controller" method:"post" summary:"Merchant Member Logout"`
}

type LogoutRes struct {
}

type PasswordResetReq struct {
	g.Meta      `path:"/passwordReset" tags:"Merchant-Member-Profile-Controller" method:"post" summary:"Merchant Member Reset Password"`
	OldPassword string `p:"oldPassword" dc:"OldPassword" v:"required"`
	NewPassword string `p:"newPassword" dc:"NewPassword" v:"required"`
}

type PasswordResetRes struct {
}
