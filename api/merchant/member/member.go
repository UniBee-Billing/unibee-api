package member

import (
	"unibee/api/bean"
	"unibee/api/bean/detail"

	"github.com/gogf/gf/v2/frame/g"
)

type ProfileReq struct {
	g.Meta `path:"/profile" tags:"Admin Member" method:"get" summary:"Get Member Profile"`
}

type ProfileRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Member Object"`
}

type UpdateReq struct {
	g.Meta    `path:"/update" tags:"Admin Member" method:"post" summary:"Update Member Profile"`
	FirstName string                  `json:"firstName"     description:"The firstName of member"`
	LastName  string                  `json:"lastName"      description:"The lastName of member"`
	Mobile    string                  `json:"mobile"     description:"mobile"`
	Metadata  *map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

type UpdateRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Member Object"`
}

type UpdateOAuthReq struct {
	g.Meta `path:"/update_oauth" tags:"Admin Member" method:"post" summary:"Update Member OAuth Account" dc:"Pass OAuth token in request header (Auth.js JWT). Headers: X-Auth-JS-Token | X-Auth-Token | X-OAuth-Token"`
}

type UpdateOAuthRes struct {
}

type ClearOAuthReq struct {
	g.Meta   `path:"/clear_oauth" tags:"Admin Member" method:"post" summary:"Clear Member OAuth Account"`
	Provider string `json:"provider" dc:"OAuth provider" v:"required"`
}

type ClearOAuthRes struct {
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
	g.Meta          `path:"/list" tags:"Admin Member" method:"get,post" summary:"Get Member List"`
	SearchKey       string   `json:"searchKey" dc:"Search Key, FirstName,LastName or Email"  `
	Email           string   `json:"email" dc:"Search Filter Email" `
	RoleIds         []uint64 `json:"roleIds" description:"The member roleId if specified'"`
	Page            int      `json:"page"  description:"Page, Start With 0" `
	Count           int      `json:"count"  description:"Count Of Page"`
	CreateTimeStart int64    `json:"createTimeStart" dc:"CreateTimeStart，UTC timestamp，seconds" `
	CreateTimeEnd   int64    `json:"createTimeEnd" dc:"CreateTimeEnd，UTC timestamp，seconds" `
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
	ReturnUrl string   `json:"returnUrl"    description:"Return url of member"`
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

type GetTotpKeyReq struct {
	g.Meta   `path:"/get_totp_key" tags:"Admin Member" method:"post" summary:"Admin Member Get 2FA Key"`
	TotpType int `json:"totpType"   description:"1-General, Google Authenticator | 2-Microsoft Authenticator | 3-Authy | 4-1Password | 5-LastPass | 6-FreeOTP | 7-Other TOTP"`
}

type GetTotpKeyRes struct {
	TotpKey        string `json:"totpKey" dc:"TotpKey, used on next confirm step"`
	TotpResumeCode string `json:"totpResumeCode" dc:"TotpResumeCode"`
	TotpUrl        string `json:"totpUrl" dc:"TotpUrl， Used to generate QR Image"`
	TotpType       int    `json:"totpType"   description:"1-General, Google Authenticator | 2-Microsoft Authenticator | 3-Authy | 4-1Password | 5-LastPass | 6-FreeOTP | 7-Other TOTP"`
}

type ConfirmTotpKeyReq struct {
	g.Meta   `path:"/confirm_totp_key" tags:"Admin Member" method:"post" summary:"Admin Member Confirm 2FA Key"`
	TotpType int    `json:"totpType"   description:"1-General, Google Authenticator | 2-Microsoft Authenticator | 3-Authy | 4-1Password | 5-LastPass | 6-FreeOTP | 7-Other TOTP"`
	TotpKey  string `json:"totpKey" description:"TotpKey"`
	TotpCode string `json:"totpCode" dc:"The totp code"`
}

type ConfirmTotpKeyRes struct {
}

type ResetTotpReq struct {
	g.Meta   `path:"/reset_totp" tags:"Admin Member" method:"post" summary:"Admin Member Reset Member 2FA Key"`
	TotpCode string `json:"totpCode" dc:"The totp code"`
}

type ResetTotpRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Member Object"`
}

type ClearTotpReq struct {
	g.Meta   `path:"/clear_member_totp" tags:"Admin Member" method:"post" summary:"Admin Owner Clear Member 2FA Key"`
	MemberId uint64 `json:"memberId"         description:"The unique id of member"`
	TotpCode string `json:"totpCode" dc:"The admin totp code"`
}

type ClearTotpRes struct {
}

type DeleteDeviceReq struct {
	g.Meta         `path:"/delete_totp_device" tags:"Admin Member" method:"post" summary:"Admin Owner Or Admin Delete Own's' 2FA Device"`
	MemberId       uint64 `json:"memberId"         description:"The unique id of member"`
	DeviceIdentity string `json:"deviceIdentity" dc:"Device Identity"`
}

type DeleteDeviceRes struct {
	DeviceList []*bean.MerchantMemberDevice `json:"deviceList" description:"The devices list'"`
}
