package _interface

import "github.com/gogf/gf/v2/net/ghttp"

type IMiddleware interface {
	ResponseHandler(r *ghttp.Request)
	PreAuth(r *ghttp.Request)
	Auth(r *ghttp.Request)
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
