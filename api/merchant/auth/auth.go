package auth

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type LoginReq struct {
	g.Meta   `path:"/sso/login" tags:"Member Authentication" method:"post" summary:"Password Login" dc:"Password login"`
	Email    string `json:"email" dc:"The merchant member email address" v:"required"`
	Password string `json:"password" dc:"The merchant member password" v:"required"`
}

type LoginRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Merchant Member Object"`
	Token          string                       `json:"token" dc:"Access token of admin portal"`
}

type LoginOtpReq struct {
	g.Meta `path:"/sso/loginOTP" tags:"Member Authentication" method:"post" summary:"OTP Login" dc:"Send email to member with OTP code"`
	Email  string `json:"email" dc:"The merchant member email address" v:"required"`
}

type LoginOtpRes struct {
}

type LoginOtpVerifyReq struct {
	g.Meta           `path:"/sso/loginOTPVerify" tags:"Member Authentication" method:"post" summary:"OTP Login Code Verification" dc:"OTP login for member, verify OTP code"`
	Email            string `json:"email" dc:"The merchant member email address" v:"required"`
	VerificationCode string `json:"verificationCode" dc:"OTP Code, received from email" v:"required"`
}

type LoginOtpVerifyRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Merchant Member Object"`
	Token          string                       `json:"token" dc:"Access token of admin portal"`
}

type PasswordForgetOtpReq struct {
	g.Meta `path:"/sso/passwordForgetOTP" tags:"Member Authentication" method:"post" summary:"OTP Password Forget" dc:"Send email to member with OTP code"`
	Email  string `json:"email" dc:"The merchant member email address" v:"required"`
}

type PasswordForgetOtpRes struct {
}

type PasswordForgetOtpVerifyReq struct {
	g.Meta           `path:"/sso/passwordForgetOTPVerify" tags:"Member Authentication" method:"post" summary:"OTP Password Forget Code Verification" dc:"Password forget OTP process, verify OTP code"`
	Email            string `json:"email" dc:"The merchant member email address" v:"required"`
	VerificationCode string `json:"verificationCode" dc:"OTP Code, received from email" v:"required"`
	NewPassword      string `json:"newPassword" dc:"The new password" v:"required"`
}

type PasswordForgetOtpVerifyRes struct {
}

type RegisterReq struct {
	g.Meta    `path:"/sso/register" tags:"Member Authentication" method:"post" summary:"Register" dc:"Register with owner permission, send email with OTP code"`
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
	g.Meta           `path:"/sso/registerVerify" tags:"Member Authentication" method:"post" summary:"Register Verify" dc:"Merchant Register, verify OTP code "`
	Email            string `json:"email" dc:"The merchant member email address" v:"required"`
	VerificationCode string `json:"verificationCode" dc:"OTP Code, received from email" v:"required"`
}

type RegisterVerifyRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Merchant Member Object"`
}
