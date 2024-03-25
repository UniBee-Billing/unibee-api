package payment

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/logic/payment/service"
)

func LinkEntry(r *ghttp.Request) {
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
