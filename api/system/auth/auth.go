package auth

import (
	"github.com/gogf/gf/v2/frame/g"
)

type TokenGeneratorReq struct {
	g.Meta        `path:"/token_generator" tags:"System-Auth" method:"post" summary:"TokenGenerator"`
	PortalType    uint64 `json:"portalType" dc:"0-Admin Portal, 1-User Portal, Default 0" default:"0"`
	MerchantId    uint64 `json:"merchantId" default:"15621"`
	Email         string `json:"email" v:"required" default:""`
	Env           string `json:"env" dc:"default daily" default:"daily"`
	RedisDatabase *int   `json:"redisDatabase" dc:"default 1" default:"1"`
}
type TokenGeneratorRes struct {
	Token string `json:"token"`
}
