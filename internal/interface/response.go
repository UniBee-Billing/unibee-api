package _interface

import (
	"github.com/gogf/gf/v2/encoding/gjson"
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
		RequestId: BizCtx().Get(r.Context()).RequestId,
	})
}

func openApiJson(r *ghttp.Request, code int, message string, data ...interface{}) {
	var responseData *gjson.Json
	if len(data) > 0 {
		responseData = gjson.New(data[0])
	} else {
		responseData = gjson.New(nil)
	}
	_ = responseData.Set("code", code)
	_ = responseData.Set("message", message)
	_ = responseData.Set("requestId", BizCtx().Get(r.Context()).RequestId)
	r.Response.WriteJson(responseData)
}

func SuccessWithMessageJsonExit(r *ghttp.Request, message string, data ...interface{}) {
	JsonExit(r, 200, "success", data)
}

func SuccessJsonExit(r *ghttp.Request, data ...interface{}) {
	JsonExit(r, 200, "success", data)
}

func FailureJsonExit(r *ghttp.Request, message string) {
	JsonExit(r, 400, message, nil)
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

// JsonRedirectExit 返回标准JSON数据引导客户端跳转，并退出当前HTTP执行函数。
func JsonRedirectExit(r *ghttp.Request, code int, message, redirect string, data ...interface{}) {
	JsonRedirect(r, code, message, redirect, data...)
	r.Exit()
}
