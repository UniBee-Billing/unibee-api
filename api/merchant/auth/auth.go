package auth

import (
	entity "unibee/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/frame/g"
)

type LoginReq struct {
	g.Meta   `path:"/sso/login" tags:"Merchant-Auth" method:"post" summary:"Login"`
	Email    string `json:"email" dc:"Email" v:"required"`
	Password string `json:"password" dc:"Password" v:"required"`
}

type LoginRes struct {
	MerchantMember *entity.MerchantMember `json:"merchantMember" dc:"Merchant Member"`
	Token          string                 `json:"token" dc:"Token"`
}

type LoginOtpReq struct {
	g.Meta `path:"/sso/loginOTP" tags:"Merchant-Auth" method:"post" summary:"Login OTP"`
	Email  string `json:"email" dc:"Email" v:"required"`
}

type LoginOtpRes struct {
}

type LoginOtpVerifyReq struct {
	g.Meta           `path:"/sso/loginOTPVerify" tags:"Merchant-Auth" method:"post" summary:"Merchant User OTP Login Verify"`
	Email            string `json:"email" dc:"Email" v:"required"`
	VerificationCode string `json:"verificationCode" dc:"VerificationCode" v:"required"`
}

type LoginOtpVerifyRes struct {
	MerchantMember *entity.MerchantMember `json:"merchantMember" dc:"Merchant Member"`
	Token          string                 `json:"token" dc:"Token"`
}

type PasswordForgetOtpReq struct {
	g.Meta `path:"/sso/passwordForgetOTP" tags:"Merchant-Auth" method:"post" summary:"Merchant Password Forget OTP"`
	Email  string `json:"email" dc:"email" v:"required"`
}

type PasswordForgetOtpRes struct {
}

type PasswordForgetOtpVerifyReq struct {
	g.Meta           `path:"/sso/passwordForgetOTPVerify" tags:"Merchant-Auth" method:"post" summary:"Merchant Password Forget OTP Verify"`
	Email            string `json:"email" dc:"email" v:"required"`
	VerificationCode string `json:"verificationCode" dc:"verificationCode" v:"required"`
	NewPassword      string `json:"newPassword" dc:"NewPassword" v:"required"`
}

type PasswordForgetOtpVerifyRes struct {
}

type RegisterReq struct {
	g.Meta    `path:"/sso/register" tags:"Merchant-Auth" method:"post" summary:"Merchant Register"`
	FirstName string `json:"firstName" dc:"First Name" v:"required"`
	LastName  string `json:"lastName" dc:"Last Name" v:"required"`
	Email     string `json:"email" dc:"Email" v:"required"`
	Password  string `json:"password" dc:"Password" v:"required"`
	Phone     string `json:"phone" dc:"Phone"`
	UserName  string `json:"userName" dc:"UserName"`
}
type RegisterRes struct {
}

type RegisterVerifyReq struct {
	g.Meta           `path:"/sso/registerVerify" tags:"Merchant-Auth" method:"post" summary:"Merchant Register Verify"`
	Email            string `json:"email" dc:"Email" v:"required"`
	VerificationCode string `json:"verificationCode" dc:"VerificationCode" v:"required"`
}

type RegisterVerifyRes struct {
	MerchantMember *entity.MerchantMember `json:"merchantMember" dc:"MerchantMember"`
}
