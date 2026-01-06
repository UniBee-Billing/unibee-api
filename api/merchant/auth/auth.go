package auth

import (
	"unibee/api/bean/detail"

	"github.com/gogf/gf/v2/frame/g"
)

type LoginReq struct {
	g.Meta     `path:"/sso/login" tags:"Member Authentication" method:"post" summary:"Password Login" dc:"Password login"`
	Email      string `json:"email" dc:"The merchant member email address" v:"required"`
	Password   string `json:"password" dc:"The merchant member password" v:"required"`
	Provider   string `json:"provider" dc:"Connect to OAuth provider"`
	ProviderId string `json:"providerId" dc:"Connect to OAuth ProviderId"`
	TotpCode   string `json:"totpCode" dc:"The totp code"`
}

type LoginRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Merchant Member Object"`
	Token          string                       `json:"token" dc:"Access token of admin portal"`
}

type LoginOAuthReq struct {
	g.Meta   `path:"/sso/oauth/login" tags:"Member Authentication" method:"post" summary:"OAuth Login" dc:"OAuth login. Pass OAuth token in request header (Auth.js JWT). Headers: X-Auth-JS-Token | X-Auth-Token | X-OAuth-Token"`
	Email    string `json:"email" dc:"The merchant member email address" v:"required"`
	TotpCode string `json:"totpCode" dc:"The totp code"`
}

type LoginOAuthRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Merchant Member Object"`
	Token          string                       `json:"token" dc:"Access token of admin portal"`
}

type SessionReq struct {
	g.Meta  `path:"/session_login" tags:"Member Authentication" method:"post" summary:"Session Login" dc:"Session login"`
	Session string `json:"session" dc:"The session" v:"required"`
}

type SessionRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Merchant Member Object"`
	Token          string                       `json:"token" dc:"Access token of admin portal"`
	ReturnUrl      string                       `json:"return_url" dc:"Return URL"`
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

type PasswordSetupOtpReq struct {
	g.Meta      `path:"/sso/passwordSetup" tags:"Member Authentication" method:"post" summary:"Password Setup" dc:"Member Password Setup"`
	Email       string `json:"email" dc:"The merchant member email address" v:"required"`
	SetupToken  string `json:"setupToken" dc:"The merchant member password setup token" v:"required"`
	NewPassword string `json:"newPassword" dc:"The new password" v:"required"`
}

type PasswordSetupOtpRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Merchant Member Object"`
	Token          string                       `json:"token" dc:"Access token of admin portal"`
}

type SetupOAuthReq struct {
	g.Meta      `path:"/sso/oauth/setup" tags:"Member Authentication" method:"post" summary:"OAuth Setup" dc:"Member OAuth Setup. Pass OAuth token in request header (Auth.js JWT). Headers: X-Auth-JS-Token | X-Auth-Token | X-OAuth-Token"`
	Email       string `json:"email" dc:"The merchant member email address" v:"required"`
	SetupToken  string `json:"setupToken" dc:"The merchant member password setup token" v:"required"`
	NewPassword string `json:"newPassword" dc:"The new password"`
}

type SetupOAuthRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Merchant Member Object"`
	Token          string                       `json:"token" dc:"Access token of admin portal"`
}

type PasswordForgetOtpVerifyReq struct {
	g.Meta           `path:"/sso/passwordForgetOTPVerify" tags:"Member Authentication" method:"post" summary:"OTP Password Forget Code Verification" dc:"Password forget OTP process, verify OTP code"`
	Email            string `json:"email" dc:"The merchant member email address" v:"required"`
	VerificationCode string `json:"verificationCode" dc:"OTP Code, received from email" v:"required"`
	NewPassword      string `json:"newPassword" dc:"The new password" v:"required"`
}

type PasswordForgetOtpVerifyRes struct {
}

type PasswordForgetTotpVerifyReq struct {
	g.Meta      `path:"/sso/passwordForgetTotpVerify" tags:"Member Authentication" method:"post" summary:"2FA Password Forget Code Verification" dc:"Password forget 2FA process, verify 2FA code"`
	Email       string `json:"email" dc:"The merchant member email address" v:"required"`
	TotpCode    string `json:"totpCode" dc:"The totp code" v:"required"`
	NewPassword string `json:"newPassword" dc:"The new password" v:"required"`
}

type PasswordForgetTotpVerifyRes struct {
}

type RegisterReq struct {
	g.Meta      `path:"/sso/register" tags:"Member Authentication" method:"post" summary:"Register" dc:"Register with owner permission, send email with OTP code"`
	FirstName   string                 `json:"firstName" dc:"The merchant owner's first name" v:"required"`
	LastName    string                 `json:"lastName" dc:"The merchant owner's last name" v:"required"`
	Email       string                 `json:"email" dc:"The merchant owner's email address" v:"required"`
	Password    string                 `json:"password" dc:"The owner's password" v:"required"`
	Phone       string                 `json:"phone" dc:"The owner's Phone"`
	UserName    string                 `json:"userName" dc:"The owner's UserName"`
	CountryCode string                 `json:"countryCode" dc:"Country Code"`
	CountryName string                 `json:"countryName" dc:"Country Name"`
	CompanyName string                 `json:"companyName" dc:"Company Name"`
	Metadata    map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}
type RegisterRes struct {
}

type RegisterEmailCheckReq struct {
	g.Meta `path:"/sso/register_email_check" tags:"Member Authentication" method:"post" summary:"Register Email Check" dc:"Check Register Email"`
	Email  string `json:"email" dc:"The merchant owner's email address" v:"required"`
}
type RegisterEmailCheckRes struct {
	Valid bool `json:"valid" dc:"Valid"`
}

type RegisterVerifyReq struct {
	g.Meta           `path:"/sso/registerVerify" tags:"Member Authentication" method:"post" summary:"Register Verify" dc:"Merchant Register, verify OTP code "`
	Email            string `json:"email" dc:"The merchant member email address" v:"required"`
	VerificationCode string `json:"verificationCode" dc:"OTP Code, received from email" v:"required"`
}

type RegisterVerifyRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Merchant Member Object"`
	Token          string                       `json:"token" dc:"Access token of admin portal"`
}

type OauthMembersReq struct {
	g.Meta `path:"/sso/oauth/members" tags:"Member Authentication" method:"get" summary:"Get Oauth Members" dc:"Merchant Get Oauth Connected Members. Pass OAuth token in request header (Auth.js JWT). Headers: X-Auth-JS-Token | X-Auth-Token | X-OAuth-Token"`
}

type OauthMembersRes struct {
	MerchantMembers []*detail.MerchantMemberDetail `json:"merchantMembers" dc:"Merchant Member Object List"`
}

type OauthGithubReq struct {
	g.Meta     `path:"/sso/oauth/github" tags:"Member Authentication" method:"get" summary:"Get Oauth Github" dc:"Merchant Oauth Github"`
	GithubCode string `json:"githubCode" dc:"Github Code"`
}

type OauthGithubRes struct {
	//GithubAccessToken string          `json:"githubAccessToken" dc:"Github Access Token"`
	UserInfo *GitHubUserInfo `json:"userInfo" dc:"GitHub User Information"`
	Sign     string          `json:"sign" dc:"Sign"`
}

type OauthGoogleReq struct {
	g.Meta      `path:"/sso/oauth/google" tags:"Member Authentication" method:"get" summary:"Get Oauth Google" dc:"Merchant Oauth Google"`
	GoogleCode  string `json:"googleCode" dc:"Google Code"`
	RedirectUri string `json:"redirectUri" dc:"The Google Redirect Uri"`
}

type OauthGoogleRes struct {
	UserInfo *GoogleUserInfo `json:"userInfo" dc:"Google User Information"`
	Sign     string          `json:"sign" dc:"Sign"`
}

type GoogleUserInfo struct {
	ID            string `json:"id" dc:"Google User ID"`
	Email         string `json:"email" dc:"User Email"`
	VerifiedEmail bool   `json:"verifiedEmail" dc:"Email Verified Status"`
	Name          string `json:"name" dc:"User Full Name"`
	GivenName     string `json:"givenName" dc:"User First Name"`
	FamilyName    string `json:"familyName" dc:"User Last Name"`
	Picture       string `json:"picture" dc:"User Profile Picture"`
	Locale        string `json:"locale" dc:"User Locale"`
}

type GitHubUserInfo struct {
	ID        int    `json:"id" dc:"GitHub User ID"`
	Login     string `json:"login" dc:"GitHub Username"`
	Email     string `json:"email" dc:"User Email"`
	Name      string `json:"name" dc:"User Full Name"`
	AvatarURL string `json:"avatarUrl" dc:"User Avatar URL"`
	Company   string `json:"company" dc:"User Company"`
	Location  string `json:"location" dc:"User Location"`
	Bio       string `json:"bio" dc:"User Bio"`
}

type RegisterOAuthReq struct {
	g.Meta      `path:"/sso/oauth/register" tags:"Member Authentication" method:"post" summary:"Register OAuth" dc:"Merchant OAuth Register. Pass OAuth token in request header (Auth.js JWT). Headers: X-Auth-JS-Token | X-Auth-Token | X-OAuth-Token"`
	Email       string                 `json:"email" dc:"The merchant member email address" v:"required"`
	Password    string                 `json:"password" dc:"The owner's password"`
	FirstName   string                 `json:"firstName" dc:"The merchant owner's first name"`
	LastName    string                 `json:"lastName" dc:"The merchant owner's last name"`
	Phone       string                 `json:"phone" dc:"The owner's Phone"`
	UserName    string                 `json:"userName" dc:"The owner's UserName"`
	CountryCode string                 `json:"countryCode" dc:"Country Code"`
	CountryName string                 `json:"countryName" dc:"Country Name"`
	CompanyName string                 `json:"companyName" dc:"Company Name"`
	Metadata    map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

type RegisterOAuthRes struct {
	MerchantMember *detail.MerchantMemberDetail `json:"merchantMember" dc:"Merchant Member Object"`
	Token          string                       `json:"token" dc:"Access token of admin portal"`
}

type ClearTotpReq struct {
	g.Meta         `path:"/sso/clear_totp" tags:"Member Authentication" method:"post" summary:"Admin Member Clear Member 2FA Key With ResumeCode"`
	Email          string `json:"email" dc:"The merchant member email address" v:"required"`
	TotpResumeCode string `json:"totpResumeCode" dc:"TOTP Resume Code" v:"required"`
}

type ClearTotpRes struct {
}
