package auth

import (
	entity "unibee/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/frame/g"
)

type LoginReq struct {
	g.Meta   `path:"/sso/login" tags:"Merchant-Auth-Controller" method:"post" summary:"Login"`
	Email    string `p:"email" dc:"Email" v:"required"`
	Password string `p:"password" dc:"Password" v:"required"`
}

type LoginRes struct {
	MerchantUser *entity.MerchantUserAccount `p:"merchantUser" dc:"Merchant User"`
	Token        string                      `p:"token" dc:"Token"`
}

type LoginOtpReq struct {
	g.Meta `path:"/sso/loginOTP" tags:"Merchant-Auth-Controller" method:"post" summary:"Login OTP"`
	Email  string `p:"email" dc:"Email" v:"required"`
}

type LoginOtpRes struct {
}

type LoginOtpVerifyReq struct {
	g.Meta           `path:"/sso/loginOTPVerify" tags:"Merchant-Auth-Controller" method:"post" summary:"Merchant User OTP Login Verify"`
	Email            string `p:"email" dc:"Email" v:"required"`
	VerificationCode string `p:"verificationCode" dc:"VerificationCode" v:"required"`
}

type LoginOtpVerifyRes struct {
	MerchantUser *entity.MerchantUserAccount `p:"merchantUser" dc:"Merchant User"`
	Token        string                      `p:"token" dc:"Token"`
}

type PasswordForgetOtpReq struct {
	g.Meta `path:"/sso/passwordForgetOTP" tags:"Merchant-Auth-Controller" method:"post" summary:"Merchant Password Forget OTP"`
	Email  string `p:"email" dc:"email" v:"required"`
}

type PasswordForgetOtpRes struct {
}

type PasswordForgetOtpVerifyReq struct {
	g.Meta           `path:"/sso/passwordForgetOTPVerify" tags:"Merchant-Auth-Controller" method:"post" summary:"Merchant Password Forget OTP Verify"`
	Email            string `p:"email" dc:"email" v:"required"`
	VerificationCode string `p:"verificationCode" dc:"verificationCode" v:"required"`
	NewPassword      string `p:"newPassword" dc:"NewPassword" v:"required"`
}

type PasswordForgetOtpVerifyRes struct {
}

type RegisterReq struct {
	g.Meta    `path:"/sso/register" tags:"Merchant-Auth-Controller" method:"post" summary:"Merchant User Register"`
	FirstName string `p:"firstName" dc:"First Name" v:"required"`
	LastName  string `p:"lastName" dc:"Last Name" v:"required"`
	Email     string `p:"email" dc:"Email" v:"required"`
	Password  string `p:"password" dc:"Password" v:"required"`
	Phone     string `p:"phone" dc:"Phone"`
	UserName  string `p:"userName" dc:"UserName"`
}
type RegisterRes struct {
}

type RegisterVerifyReq struct {
	g.Meta           `path:"/sso/registerVerify" tags:"Merchant-Auth-Controller" method:"post" summary:"Merchant Register Verify"`
	Email            string `p:"email" dc:"Email" v:"required"`
	VerificationCode string `p:"verificationCode" dc:"VerificationCode" v:"required"`
}

type RegisterVerifyRes struct {
	MerchantUser *entity.MerchantUserAccount `p:"merchantUser" dc:"MerchantUser"`
}
