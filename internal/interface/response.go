package _interface

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type JsonRes struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Redirect  string      `json:"redirect"`
	RequestId string      `json:"requestId"`
}

func portalJson(r *ghttp.Request, code int, message string, data ...interface{}) {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	} else {
		responseData = g.Map{}
	}
	r.Response.WriteJson(JsonRes{
		Code:      code,
		Message:   message,
		Data:      responseData,
		RequestId: Context().Get(r.Context()).RequestId,
	})
}

func JsonExit(r *ghttp.Request, code int, message string, data ...interface{}) {
	portalJson(r, code, message, data...)
	r.Exit()
}

func OpenApiJsonExit(r *ghttp.Request, code int, message string, data ...interface{}) {
	portalJson(r, code, message, data...)
	r.Exit()
}

func JsonRedirect(r *ghttp.Request, code int, message, redirect string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	r.Response.WriteJson(JsonRes{
		Code:     code,
		Message:  message,
		Data:     responseData,
		Redirect: redirect,
	})
}

func JsonRedirectExit(r *ghttp.Request, code int, message, redirect string, data ...interface{}) {
	JsonRedirect(r, code, message, redirect, data...)
	r.Exit()
}
