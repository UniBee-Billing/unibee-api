package auth

import (
	entity "go-oversea-pay/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/frame/g"
)

type LoginReq struct {
	g.Meta   `path:"/sso/login" tags:"User-Auth-Controller" method:"post" summary:"1.1 用户登录"`
	Email    string `p:"email" dc:"email" v:"required"`
	Password string `p:"password" dc:"password" v:"required"`
}

type LoginRes struct {
	User  *entity.UserAccount `p:"user" dc:"user"`
	Token string              `p:"token" dc:"token string"`
}

type LoginOtpReq struct {
	g.Meta `path:"/sso/loginOTP" tags:"User-Auth-Controller" method:"post" summary:"1.1 用户OTP登录"`
	Email  string `p:"email" dc:"email" v:"required"`
}

type LoginOtpRes struct {
}

type LoginOtpVerifyReq struct {
	g.Meta           `path:"/sso/loginOTPVerify" tags:"User-Auth-Controller" method:"post" summary:"1.1 用户OTP登录"`
	Email            string `p:"email" dc:"email" v:"required"`
	VerificationCode string `p:"verificationCode" dc:"verificationCode" v:"required"`
}

type LoginOtpVerifyRes struct {
	User  *entity.UserAccount `p:"user" dc:"user"`
	Token string              `p:"token" dc:"token string"`
}

type LogoutReq struct {
	g.Meta `path:"/sso/logout" tags:"User-Auth-Controller" method:"post" summary:"1.1 user logout"`
}

type LogoutRes struct {
}

type RegisterReq struct {
	g.Meta      `path:"/sso/register" tags:"User-Auth-Controller" method:"post" summary:"1.1 用户注册"`
	FirstName   string `p:"firstName" dc:"first name" v:"required"`
	LastName    string `p:"lastName" dc:"last name" v:"required"`
	Email       string `p:"email" dc:"email" v:"required"`
	Password    string `p:"password" dc:"password" v:"required"`
	Phone       string `p:"phone" dc:"phone" `
	Address     string `p:"address" dc:"adderss"`
	CountryCode string `p:"countryCode" dc:"countryCode"`
	CountryName string `p:"countryName" dc:"countryName"`
	UserName    string `p:"userName" dc:"userName" v:"required"`
}
type RegisterRes struct {
	// User *entity.UserAccount `p:"user" dc:"user"`
}

type RegisterVerifyReq struct {
	g.Meta           `path:"/sso/registerVerify" tags:"User-Auth-Controller" method:"post" summary:"1.2 用户注册(verify email)"`
	Email            string `p:"email" dc:"email" v:"required"`
	VerificationCode string `p:"verificationCode" dc:"verification code" v:"required"`
}

type RegisterVerifyRes struct {
	User *entity.UserAccount `p:"user" dc:"user"`
}
