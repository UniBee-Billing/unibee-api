package model

import (
	entity "unibee/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type Context struct {
	Session        *ghttp.Session
	MerchantId     uint64
	User           *ContextUser
	MerchantMember *ContextMerchantMember
	RequestId      string
	Data           g.Map
	OpenApiConfig  *entity.OpenApiConfig
	IsOpenApiCall  bool
}

type ContextUser struct {
	Id         uint64
	MerchantId uint64
	Token      string
	Email      string
}

type ContextMerchantMember struct {
	Id         uint64
	MerchantId uint64
	Token      string
	Email      string
	IsOwner    bool
}
