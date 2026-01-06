package model

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/golang-jwt/jwt/v5"
	entity "unibee/internal/model/entity/default"
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

type OpenApiConfig struct {
	Id                      uint64 `json:"id"                      description:""`                         //
	Qps                     int    `json:"qps"                     description:"total qps control"`        // total qps control
	MerchantId              uint64 `json:"merchantId"              description:"merchant id"`              // merchant id
	Hmac                    string `json:"hmac"                    description:"webhook hmac key"`         // webhook hmac key
	Callback                string `json:"callback"                description:"callback url"`             // callback url
	ApiKey                  string `json:"apiKey"                  description:"api key"`                  // api key
	Token                   string `json:"token"                   description:"api token"`                // api token
	IsDeleted               int    `json:"isDeleted"               description:"0-UnDeleted，1-Deleted"`    // 0-UnDeleted，1-Deleted
	ValidIps                string `json:"validIps"                description:""`                         //
	GatewayCallbackResponse string `json:"gatewayCallbackResponse" description:"callback return response"` // callback return response
	CompanyId               int64  `json:"companyId"               description:"company id"`               // company id
}

type Context struct {
	Session           *ghttp.Session
	MerchantId        uint64
	User              *ContextUser
	MerchantMember    *ContextMerchantMember
	RequestId         string
	Data              g.Map
	OpenApiConfig     *OpenApiConfig
	OpenApiKey        string
	IsOpenApiCall     bool
	IsAdminPortalCall bool
	Language          string
	UserAgent         string
	UserAgentIdentity string
	RemoteIP          string
	ClientIdentity    string
	Authorization     string
	TokenString       string
	Token             *TokenClaims
	PreloadData       *PreloadData
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

type PreloadData struct {
	Plans                   map[uint64]*entity.Plan
	Products                map[uint64]*entity.Product
	Payments                map[string]*entity.Payment
	Refunds                 map[string]*entity.Refund
	Invoices                map[string]*entity.Invoice
	Gateways                map[uint64]*entity.MerchantGateway
	Users                   map[uint64]*entity.UserAccount
	Subscriptions           map[string]*entity.Subscription
	PromoCreditTransactions map[string]*entity.CreditTransaction
	Discounts               map[string]*entity.MerchantDiscountCode
}
