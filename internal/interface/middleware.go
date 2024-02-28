package _interface

import "github.com/gogf/gf/v2/net/ghttp"

type IMiddleware interface {
	CORS(r *ghttp.Request)
	ResponseHandler(r *ghttp.Request)
	OpenApiDetach(r *ghttp.Request)
	TokenAuth(r *ghttp.Request)
	UserPortalPreAuth(r *ghttp.Request)
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
