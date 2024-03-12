package member

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type ProfileReq struct {
	g.Meta `path:"/profile" tags:"Member" method:"get" summary:"Get Merchant Member Profile"`
}

type ProfileRes struct {
	MerchantMember *bean.MerchantMemberSimplify `json:"merchantMember" dc:"Merchant Member"`
}

type LogoutReq struct {
	g.Meta `path:"/logout" tags:"Member" method:"post" summary:"Merchant Member Logout"`
}

type LogoutRes struct {
}

type PasswordResetReq struct {
	g.Meta      `path:"/passwordReset" tags:"Member" method:"post" summary:"Merchant Member Reset Password"`
	OldPassword string `json:"oldPassword" dc:"OldPassword" v:"required"`
	NewPassword string `json:"newPassword" dc:"NewPassword" v:"required"`
}

type PasswordResetRes struct {
}
