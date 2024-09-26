package _interface

import "github.com/gogf/gf/v2/net/ghttp"

type IMiddleware interface {
	CORS(r *ghttp.Request)
	ResponseHandler(r *ghttp.Request)
	UserPortalApiHandler(r *ghttp.Request)
	MerchantHandler(r *ghttp.Request)
	UserPortalMerchantRouterHandler(r *ghttp.Request)
}

var localMiddleware IMiddleware

func Middleware() IMiddleware {
	if localMiddleware == nil {
		panic("implement not found for interface IMiddleware, forgot register?")
	}
	return localMiddleware
}

func RegisterMiddleware(i IMiddleware) {
	localMiddleware = i
}
