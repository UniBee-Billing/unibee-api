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

type ListReq struct {
	g.Meta `path:"/list" tags:"Member" method:"get" summary:"Get Merchant Member List"`
}

type ListRes struct {
	MerchantMembers []*bean.MerchantMemberSimplify `json:"merchantMembers" dc:"Merchant Members"`
}

type UpdateMemberRoleReq struct {
	g.Meta   `path:"/update_member_role" tags:"Member" method:"post" summary:"Update Member Role"`
	MemberId uint64 `json:"memberId"         description:"MemberId"`
	Role     string `json:"role"         description:"Role"`
}

type UpdateMemberRoleRes struct {
}

type NewMemberReq struct {
	g.Meta    `path:"/new_member" tags:"Member" method:"post" summary:"Invite member, will send email to member email provided, member can enter admin portal by email otp login"`
	Email     string `json:"email"  v:"required"   description:"Email, email of member, member otp login needs" `
	Role      string `json:"role"    v:"required"     description:"Role, permission role of member" `
	FirstName string `json:"firstName"     description:"FirstName, member firstName, invoice needs"`
	LastName  string `json:"lastName"      description:"LastName, member firstName, invoice needs"`
}

type NewMemberRes struct {
}
