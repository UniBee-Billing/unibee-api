package auth

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type LoginReq struct {
	g.Meta   `path:"/sso/login" tags:"Auth" method:"post" summary:"Login" dc:"Password login for merchant member'"`
	Email    string `json:"email" dc:"The merchant member's email address'" v:"required"`
	Password string `json:"password" dc:"The merchant member's password'" v:"required"`
}

type LoginRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Merchant Member Object"`
	Token          string                       `json:"token" dc:"Access token of admin portal"`
}

type LoginOtpReq struct {
	g.Meta `path:"/sso/loginOTP" tags:"Auth" method:"post" summary:"LoginOTP" dc:"OTP login for merchant member, send email to member's email address with OTP code'"`
	Email  string `json:"email" dc:"The merchant member's email address" v:"required"`
}

type LoginOtpRes struct {
}

type LoginOtpVerifyReq struct {
	g.Meta           `path:"/sso/loginOTPVerify" tags:"Auth" method:"post" summary:"LoginOTPVerify" dc:"OTP login for merchant member, verify OTP code"`
	Email            string `json:"email" dc:"The merchant member's email address" v:"required"`
	VerificationCode string `json:"verificationCode" dc:"OTP Code, received from email" v:"required"`
}

type LoginOtpVerifyRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Merchant Member Object"`
	Token          string                       `json:"token" dc:"Access token of admin portal"`
}

type PasswordForgetOtpReq struct {
	g.Meta `path:"/sso/passwordForgetOTP" tags:"Auth" method:"post" summary:"PasswordForgetOTP" dc:"Merchant member's password forget OTP process,, send email to member's email address with OTP code'"`
	Email  string `json:"email" dc:"The merchant member's email address" v:"required"`
}

type PasswordForgetOtpRes struct {
}

type PasswordForgetOtpVerifyReq struct {
	g.Meta           `path:"/sso/passwordForgetOTPVerify" tags:"Auth" method:"post" summary:"PasswordForgetOTPVerify" dc:"Merchant member's password forget OTP process, verify OTP code"`
	Email            string `json:"email" dc:"The merchant member's email address" v:"required"`
	VerificationCode string `json:"verificationCode" dc:"OTP Code, received from email" v:"required"`
	NewPassword      string `json:"newPassword" dc:"The new password" v:"required"`
}

type PasswordForgetOtpVerifyRes struct {
}

type RegisterReq struct {
	g.Meta    `path:"/sso/register" tags:"Auth" method:"post" summary:"Register" dc:"Register merchant with owner, send email to owner's email address with OTP code, only open for cloud version; the owner account will create automatic for standalone version"`
	FirstName string `json:"firstName" dc:"The merchant owner's first name" v:"required"`
	LastName  string `json:"lastName" dc:"The merchant owner's last name" v:"required"`
	Email     string `json:"email" dc:"The merchant owner's email address" v:"required"`
	Password  string `json:"password" dc:"The owner's password" v:"required"`
	Phone     string `json:"phone" dc:"The owner's Phone"`
	UserName  string `json:"userName" dc:"The owner's UserName"`
}
type RegisterRes struct {
}

type RegisterVerifyReq struct {
	g.Meta           `path:"/sso/registerVerify" tags:"Auth" method:"post" summary:"RegisterVerify" dc:"Merchant Register, verify OTP code "`
	Email            string `json:"email" dc:"The merchant member's email address" v:"required"`
	VerificationCode string `json:"verificationCode" dc:"OTP Code, received from email" v:"required"`
}

type RegisterVerifyRes struct {
	MerchantMember *bean.MerchantMemberSimplify `json:"merchantMember" dc:"Merchant Member Object"`
}
