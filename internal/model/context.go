package model

import (
	entity "unibee-api/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type Context struct {
	Session       *ghttp.Session
	MerchantId    uint64
	User          *ContextUser
	MerchantUser  *ContextMerchantUser
	RequestId     string
	Data          g.Map
	OpenApiConfig *entity.OpenApiConfig
}

type ContextUser struct {
	Id         uint64
	MerchantId uint64
	Token      string
	Email      string
}

type ContextMerchantUser struct {
	Id         uint64
	MerchantId uint64
	Token      string
	Email      string
}
