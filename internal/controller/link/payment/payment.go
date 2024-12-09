package payment

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/logic/payment/service"
)

func LinkEntry(r *ghttp.Request) {
	r.Response.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	r.Response.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT,DELETE,OPTIONS,PATCH")
	r.Response.Header().Add("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		return
	}
	paymentId := r.Get("paymentId").String()
	res := service.LinkCheck(r.Context(), paymentId, gtime.Now().Timestamp())
	if len(res.Link) > 0 {
		r.Response.RedirectTo(res.Link)
	} else if len(res.Message) > 0 {
		r.Response.Writeln(res.Message)
	} else {
		r.Response.Writeln("Server Error")
	}
}
