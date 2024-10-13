package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
	"unibee/utility"
)

func HealthCheck(r *ghttp.Request) {
	r.Response.WriteHeader(http.StatusOK)
	return
}

func Version(r *ghttp.Request) {
	r.Response.WriteHeader(http.StatusOK)
	r.Response.Write(utility.ReadBuildVersionInfo(r.Context()))
	return
}
