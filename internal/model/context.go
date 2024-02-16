package model

import (
	entity "unibee-api/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Context 请求上下文结构
type Context struct {
	Session       *ghttp.Session // 当前Session管理对象
	User          *ContextUser   // 上下文用户信息
	MerchantUser  *ContextMerchantUser
	RequestId     string
	Data          g.Map // 自定KV变量，业务模块根据需要设置，不固定
	OpenApiConfig *entity.OpenApiConfig
}

// ContextUser 请求上下文中的用户信息
type ContextUser struct {
	Id    uint64 // 用户ID
	Token string // token
	// MobilePhone string // 用户手机号
	// UserName    string // 用户名称
	// AvatarUrl   string // 用户头像
	// IsAdmin     bool   // 是否是管理员
	Email string
}

type ContextMerchantUser struct {
	Id         uint64 // 用户ID
	MerchantId uint64 // MerchantId
	Token      string // token
	// MobilePhone string // 用户手机号
	// UserName    string // 用户名称
	// AvatarUrl   string // 用户头像
	// IsAdmin     bool   // 是否是管理员
	Email string
}
