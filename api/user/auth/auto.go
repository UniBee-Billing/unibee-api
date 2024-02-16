package auth

import (
	entity "unibee-api/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/frame/g"
)

type LoginReq struct {
	g.Meta   `path:"/sso/login" tags:"User-Auth-Controller" method:"post" summary:"User Login"`
	Email    string `p:"email" dc:"email" v:"required"`
	Password string `p:"password" dc:"password" v:"required"`
}

type LoginRes struct {
	User  *entity.UserAccount `p:"user" dc:"user"`
	Token string              `p:"token" dc:"token string"`
}

type LoginOtpReq struct {
	g.Meta `path:"/sso/loginOTP" tags:"User-Auth-Controller" method:"post" summary:"User OTP Login"`
	Email  string `p:"email" dc:"email" v:"required"`
}

type LoginOtpRes struct {
}

type LoginOtpVerifyReq struct {
	g.Meta           `path:"/sso/loginOTPVerify" tags:"User-Auth-Controller" method:"post" summary:"User OTP Login"`
	Email            string `p:"email" dc:"email" v:"required"`
	VerificationCode string `p:"verificationCode" dc:"verificationCode" v:"required"`
}

type LoginOtpVerifyRes struct {
	User  *entity.UserAccount `p:"user" dc:"user"`
	Token string              `p:"token" dc:"token"`
}

type RegisterReq struct {
	g.Meta      `path:"/sso/register" tags:"User-Auth-Controller" method:"post" summary:"User Register"`
	FirstName   string `p:"firstName" dc:"First Name" v:"required"`
	LastName    string `p:"lastName" dc:"Last Name" v:"required"`
	Email       string `p:"email" dc:"Email" v:"required"`
	Password    string `p:"password" dc:"Password" v:"required"`
	Phone       string `p:"phone" dc:"Phone" `
	Address     string `p:"address" dc:"Address"`
	CountryCode string `p:"countryCode" dc:"CountryCode"`
	CountryName string `p:"countryName" dc:"CountryName"`
	UserName    string `p:"userName" dc:"UserName"`
}
type RegisterRes struct {
	// User *entity.UserAccount `p:"user" dc:"user"`
}

type RegisterVerifyReq struct {
	g.Meta           `path:"/sso/registerVerify" tags:"User-Auth-Controller" method:"post" summary:"User Register Via Email"`
	Email            string `p:"email" dc:"Email" v:"required"`
	VerificationCode string `p:"verificationCode" dc:"Verification Code" v:"required"`
}

type RegisterVerifyRes struct {
	User *entity.UserAccount `p:"user" dc:"User"`
}
