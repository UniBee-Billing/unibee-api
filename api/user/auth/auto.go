package auth

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type LoginReq struct {
	g.Meta   `path:"/sso/login" tags:"User-Auth" method:"post" summary:"User Login"`
	Email    string `json:"email" dc:"email" v:"required"`
	Password string `json:"password" dc:"password" v:"required"`
}

type LoginRes struct {
	User  *bean.UserAccountSimplify `json:"user" dc:"user"`
	Token string                    `json:"token" dc:"token string"`
}

type SessionLoginReq struct {
	g.Meta  `path:"/session_login" tags:"User-Auth" method:"post" summary:"User Portal Session Login"`
	Session string `json:"session" dc:"Session" v:"required"`
}

type SessionLoginRes struct {
	User  *bean.UserAccountSimplify `json:"user" dc:"user"`
	Token string                    `json:"token" dc:"token string"`
}

type LoginOtpReq struct {
	g.Meta `path:"/sso/loginOTP" tags:"User-Auth" method:"post" summary:"User OTP Login"`
	Email  string `json:"email" dc:"email" v:"required"`
}

type LoginOtpRes struct {
}

type LoginOtpVerifyReq struct {
	g.Meta           `path:"/sso/loginOTPVerify" tags:"User-Auth" method:"post" summary:"User OTP Login Verify"`
	Email            string `json:"email" dc:"email" v:"required"`
	VerificationCode string `json:"verificationCode" dc:"verificationCode" v:"required"`
}

type LoginOtpVerifyRes struct {
	User  *bean.UserAccountSimplify `json:"user" dc:"user"`
	Token string                    `json:"token" dc:"token"`
}

type PasswordForgetOtpReq struct {
	g.Meta `path:"/sso/passwordForgetOTP" tags:"User-Auth" method:"post" summary:"User Password Forget OTP"`
	Email  string `json:"email" dc:"email" v:"required"`
}

type PasswordForgetOtpRes struct {
}

type PasswordForgetOtpVerifyReq struct {
	g.Meta           `path:"/sso/passwordForgetOTPVerify" tags:"User-Auth" method:"post" summary:"User Password Forget OTP Verify"`
	Email            string `json:"email" dc:"email" v:"required"`
	VerificationCode string `json:"verificationCode" dc:"verificationCode" v:"required"`
	NewPassword      string `json:"newPassword" dc:"NewPassword" v:"required"`
}

type PasswordForgetOtpVerifyRes struct {
}

type RegisterReq struct {
	g.Meta      `path:"/sso/register" tags:"User-Auth" method:"post" summary:"User Register"`
	FirstName   string `json:"firstName" dc:"First Name" v:"required"`
	LastName    string `json:"lastName" dc:"Last Name" v:"required"`
	Email       string `json:"email" dc:"Email" v:"required"`
	Password    string `json:"password" dc:"Password" v:"required"`
	Phone       string `json:"phone" dc:"Phone" `
	Address     string `json:"address" dc:"Address"`
	CountryCode string `json:"countryCode" dc:"CountryCode"`
	CountryName string `json:"countryName" dc:"CountryName"`
	UserName    string `json:"userName" dc:"UserName"`
}
type RegisterRes struct {
}

type RegisterVerifyReq struct {
	g.Meta           `path:"/sso/registerVerify" tags:"User-Auth" method:"post" summary:"User Register Verify"`
	Email            string `json:"email" dc:"Email" v:"required"`
	VerificationCode string `json:"verificationCode" dc:"Verification Code" v:"required"`
}

type RegisterVerifyRes struct {
	User *bean.UserAccountSimplify `json:"user" dc:"User"`
}
