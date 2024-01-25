package auth

import (
	entity "go-oversea-pay/internal/model/entity/oversea_pay"

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
	g.Meta           `path:"/sso/loginOTPVerify" tags:"Merchant-Auth-Controller" method:"post" summary:"User OTP Login Verify Code"`
	Email            string `p:"email" dc:"Email" v:"required"`
	VerificationCode string `p:"verificationCode" dc:"VerificationCode" v:"required"`
}

type LoginOtpVerifyRes struct {
	MerchantUser *entity.MerchantUserAccount `p:"merchantUser" dc:"Merchant User"`
	Token        string                      `p:"token" dc:"Token"`
}

type RegisterReq struct {
	g.Meta     `path:"/sso/register" tags:"Merchant-Auth-Controller" method:"post" summary:"User Register"`
	FirstName  string `p:"firstName" dc:"First Name" v:"required"`
	LastName   string `p:"lastName" dc:"Last Name" v:"required"`
	Email      string `p:"email" dc:"Email" v:"required"`
	Password   string `p:"password" dc:"Password" v:"required"`
	Phone      string `p:"phone" dc:"Phone"`
	MerchantId uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
	UserName   string `p:"userName" dc:"UserName"`
	// Address   string `p:"address" dc:"adderss"`
}
type RegisterRes struct {
	// User *entity.MerchantUserAccount `p:"user" dc:"user"`
}

type RegisterVerifyReq struct {
	g.Meta           `path:"/sso/registerVerify" tags:"Merchant-Auth-Controller" method:"post" summary:"Verify Email"`
	Email            string `p:"email" dc:"Email" v:"required"`
	VerificationCode string `p:"verificationCode" dc:"VerificationCode" v:"required"`
}

// NO, after successful signup, res should be empty, front-end should be redirectd to /login
type RegisterVerifyRes struct {
	MerchantUser *entity.MerchantUserAccount `p:"merchantUser" dc:"MerchantUser"`
}
