package model

import (
	"github.com/golang-jwt/jwt/v5"
	entity "unibee/internal/model/entity/default"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type TokenType string

type TokenClaims struct {
	TokenType     TokenType `json:"tokenType"`
	Id            uint64    `json:"id"`
	Email         string    `json:"email"`
	MerchantId    uint64    `json:"merchantId"`
	PermissionKey string    `json:"permissionKey"`
	Lang          string    `json:"lang"`
	jwt.RegisteredClaims
}

type Context struct {
	Session        *ghttp.Session
	MerchantId     uint64
	User           *ContextUser
	MerchantMember *ContextMerchantMember
	RequestId      string
	Data           g.Map
	OpenApiConfig  *entity.OpenApiConfig
	OpenApiKey     string
	IsOpenApiCall  bool
	Language       string
	UserAgent      string
	Authorization  string
	TokenString    string
	Token          *TokenClaims
}

type ContextUser struct {
	Id         uint64
	MerchantId uint64
	Token      string
	Email      string
	Lang       string
}

type ContextMerchantMember struct {
	Id         uint64
	MerchantId uint64
	Token      string
	Email      string
	IsOwner    bool
}
