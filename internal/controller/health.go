package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
)

func HealthCheck(r *ghttp.Request) {
	r.Response.WriteHeader(http.StatusOK)
	return
}
