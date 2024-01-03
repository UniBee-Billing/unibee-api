package _interface

import "github.com/gogf/gf/v2/net/ghttp"

type IMiddleware interface {
	ResponseHandler(r *ghttp.Request)
	PreAuth(r *ghttp.Request)
	PreOpenApiAuth(r *ghttp.Request)
	Auth(r *ghttp.Request)
	CORS(r *ghttp.Request)
	TokenAuth(r *ghttp.Request)
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
