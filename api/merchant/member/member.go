package member

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type ProfileReq struct {
	g.Meta `path:"/profile" tags:"Admin Member" method:"get" summary:"Get Member Profile"`
}

type ProfileRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Member Object"`
}

type UpdateReq struct {
	g.Meta    `path:"/update" tags:"Admin Member" method:"post" summary:"Update Member Profile"`
	FirstName string `json:"firstName"     description:"The firstName of member"`
	LastName  string `json:"lastName"      description:"The lastName of member"`
	Mobile    string `json:"mobile"     description:"mobile"`
}

type UpdateRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Member Object"`
}

type LogoutReq struct {
	g.Meta `path:"/logout" tags:"Admin Member" method:"post" summary:"Logout"`
}

type LogoutRes struct {
}

type PasswordResetReq struct {
	g.Meta      `path:"/passwordReset" tags:"Admin Member" method:"post" summary:"Member Reset Password"`
	OldPassword string `json:"oldPassword" dc:"The old password of member" v:"required"`
	NewPassword string `json:"newPassword" dc:"The new password of member" v:"required"`
}

type PasswordResetRes struct {
}

type ListReq struct {
	g.Meta  `path:"/list" tags:"Admin Member" method:"get,post" summary:"Get Member List"`
	RoleIds []uint64 `json:"roleIds" description:"The member roleId if specified'"`
	Page    int      `json:"page"  description:"Page, Start With 0" `
	Count   int      `json:"count"  description:"Count Of Page"`
}

type ListRes struct {
	MerchantMembers []*detail.MerchantMemberDetail `json:"merchantMembers" dc:"Merchant Member Object List"`
	Total           int                            `json:"total" dc:"Total"`
}

type UpdateMemberRoleReq struct {
	g.Meta   `path:"/update_member_role" tags:"Admin Member" method:"post" summary:"Update Member Role"`
	MemberId uint64   `json:"memberId"         description:"The unique id of member"`
	RoleIds  []uint64 `json:"roleIds"         description:"The id list of role"`
}

type UpdateMemberRoleRes struct {
}

type NewMemberReq struct {
	g.Meta    `path:"/new_member" tags:"Admin Member" method:"post" summary:"Invite member" description:"Will send email to member email provided, member can enter admin portal by email otp login"`
	Email     string   `json:"email"  v:"required"   description:"The email of member" `
	RoleIds   []uint64 `json:"roleIds"    v:"required"     description:"The id list of role" `
	FirstName string   `json:"firstName"     description:"The firstName of member"`
	LastName  string   `json:"lastName"      description:"The lastName of member"`
}

type NewMemberRes struct {
}

type FrozenReq struct {
	g.Meta   `path:"/suspend_member" tags:"Admin Member" method:"post" summary:"Suspend Member"`
	MemberId uint64 `json:"memberId"         description:"The unique id of member"`
}

type FrozenRes struct {
}

type ReleaseReq struct {
	g.Meta   `path:"/resume_member" tags:"Admin Member" method:"post" summary:"Resume Member"`
	MemberId uint64 `json:"memberId"         description:"The unique id of member"`
}

type ReleaseRes struct {
}

type OperationLogListReq struct {
	g.Meta          `path:"/operation_log_list" tags:"Admin Member" method:"get" summary:"Get Member Operation Log List"`
	MemberFirstName string `json:"memberFirstName" description:"Filter Member's FirstName Default All" `
	MemberLastName  string `json:"memberLastName" description:"Filter Member's LastName, Default All" `
	MemberEmail     string `json:"memberEmail" description:"Filter Member's Email, Default All" `
	FirstName       string `json:"firstName" description:"FirstName" `
	LastName        string `json:"lastName" description:"LastName" `
	Email           string `json:"email" description:"Email" `
	SubscriptionId  string `json:"subscriptionId"     description:"subscription_id"` // subscription_id
	InvoiceId       string `json:"invoiceId"          description:"invoice id"`      // invoice id
	PlanId          uint64 `json:"planId"             description:"plan id"`         // plan id
	DiscountCode    string `json:"discountCode"       description:"discount_code"`   // discount_code
	Page            int    `json:"page"  description:"Page, Start With 0" `
	Count           int    `json:"count"  description:"Count Of Page"`
}

type OperationLogListRes struct {
	MerchantOperationLogs []*detail.MerchantOperationLogDetail `json:"merchantOperationLogs" dc:"Merchant Member Operation Log List"`
	Total                 int                                  `json:"total" dc:"Total"`
}
