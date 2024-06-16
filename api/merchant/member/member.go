package member

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type ProfileReq struct {
	g.Meta `path:"/profile" tags:"Member" method:"get" summary:"GetMemberProfile"`
}

type ProfileRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Member Object"`
}

type LogoutReq struct {
	g.Meta `path:"/logout" tags:"Member" method:"post" summary:"MemberLogout"`
}

type LogoutRes struct {
}

type PasswordResetReq struct {
	g.Meta      `path:"/passwordReset" tags:"Member" method:"post" summary:"MemberResetPassword"`
	OldPassword string `json:"oldPassword" dc:"The old password of member" v:"required"`
	NewPassword string `json:"newPassword" dc:"The new password of member" v:"required"`
}

type PasswordResetRes struct {
}

type ListReq struct {
	g.Meta `path:"/list" tags:"Member" method:"get" summary:"GetMemberList"`
}

type ListRes struct {
	MerchantMembers []*detail.MerchantMemberDetail `json:"merchantMembers" dc:"Merchant Member Object List"`
	Total           int                            `json:"total" dc:"Total"`
}

type UpdateMemberRoleReq struct {
	g.Meta   `path:"/update_member_role" tags:"Member" method:"post" summary:"UpdateMemberRole"`
	MemberId uint64   `json:"memberId"         description:"The unique id of member"`
	RoleIds  []uint64 `json:"roleIds"         description:"The id list of role"`
}

type UpdateMemberRoleRes struct {
}

type NewMemberReq struct {
	g.Meta    `path:"/new_member" tags:"Member" method:"post" summary:"Invite member" description:"Will send email to member email provided, member can enter admin portal by email otp login"`
	Email     string   `json:"email"  v:"required"   description:"The email of member" `
	RoleIds   []uint64 `json:"roleIds"    v:"required"     description:"The id list of role" `
	FirstName string   `json:"firstName"     description:"The firstName of member"`
	LastName  string   `json:"lastName"      description:"The lastName of member"`
}

type NewMemberRes struct {
}
