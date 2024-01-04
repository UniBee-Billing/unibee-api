package auth

import (
	entity "go-oversea-pay/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/frame/g"
)

type LoginReq struct {
	g.Meta   `path:"/sso/login" tags:"MerchantUser-Auth-Controller" method:"post" summary:"1.1 用户登录"`
	Email    string `p:"email" dc:"email" v:"required"`
	Password string `p:"password" dc:"password" v:"required"`
}

type LoginRes struct {
	MerchantUser  *entity.MerchantUserAccount `p:"merchantUser" dc:"merchant user"`
	Token string              `p:"token" dc:"token string"`
}

type LoginOtpReq struct {
	g.Meta `path:"/sso/loginOTP" tags:"MerchantUser-Auth-Controller" method:"post" summary:"1.1 用户OTP登录"`
	Email  string `p:"email" dc:"email" v:"required"`
}

type LoginOtpRes struct {
}

type LoginOtpVerifyReq struct {
	g.Meta           `path:"/sso/loginOTPVerify" tags:"MerchantUser-Auth-Controller" method:"post" summary:"1.1 用户OTP登录"`
	Email            string `p:"email" dc:"email" v:"required"`
	VerificationCode string `p:"verificationCode" dc:"verificationCode" v:"required"`
}

type LoginOtpVerifyRes struct {
	MerchantUser  *entity.MerchantUserAccount `p:"merchantUser" dc:"merchant user"`
	Token string              `p:"token" dc:"token string"`
}

type LogoutReq struct {
	g.Meta           `path:"/sso/logout" tags:"User-Auth-Controller" method:"post" summary:"1.1 user logout"`
}

type LogoutRes struct {
}

type RegisterReq struct {
	g.Meta    `path:"/sso/register" tags:"User-Auth-Controller" method:"post" summary:"1.1 用户注册"`
	FirstName string `p:"firstName" dc:"first name" v:"required"`
	LastName  string `p:"lastName" dc:"last name" v:"required"`
	Email     string `p:"email" dc:"email" v:"required"`
	Password  string `p:"password" dc:"password" v:"required"`
	Phone     string `p:"phone" dc:"phone"`
	MerchantId	uint64 	`p:"merchantId" dc:"merchant id" v:"required"`
	UserName	string `p:"userName" dc:"user name"`
	// Address   string `p:"address" dc:"adderss"`
}
type RegisterRes struct {
	// User *entity.MerchantUserAccount `p:"user" dc:"user"`
}

type RegisterVerifyReq struct {
	g.Meta           `path:"/sso/registerVerify" tags:"User-Auth-Controller" method:"post" summary:"1.2 用户注册(verify email)"`
	Email            string `p:"email" dc:"email" v:"required"`
	VerificationCode string `p:"verificationCode" dc:"verification code" v:"required"`
}

// NO, after successful signup, res should be empty, front-end should be redirectd to /login
type RegisterVerifyRes struct {
	MerchantUser *entity.MerchantUserAccount `p:"merchantUser" dc:"merchant user"`
}
