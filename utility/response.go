package utility

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	_interface "unibee/internal/interface"
)

// JsonRes 数据返回通用JSON数据结构
type JsonRes struct {
	Code      int         `json:"code"`     // 错误码((0:成功, 1:失败, >1:错误码))
	Message   string      `json:"message"`  // 提示信息
	Data      interface{} `json:"data"`     // 返回数据(业务接口定义具体数据结构)
	Redirect  string      `json:"redirect"` // 引导客户端跳转到指定路由
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
		RequestId: _interface.BizCtx().Get(r.Context()).RequestId,
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
	_ = responseData.Set("requestId", _interface.BizCtx().Get(r.Context()).RequestId)
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
	openApiJson(r, code, message, data...)
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
