package payment

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/internal/consts"
	"unibee/internal/query"
)

func PaymentLinkEntry(r *ghttp.Request) {
	paymentId := r.Get("paymentId").String()
	one := query.GetPaymentByPaymentId(r.Context(), paymentId)
	if one == nil {
		g.Log().Errorf(r.Context(), "PaymentLinkEntry payment not found url: %s paymentId: %s", r.GetUrl(), paymentId)
		return
	}
	if one.Status == consts.PaymentCancelled {
		r.Response.Writeln("Payment Cancelled")
	} else if one.Status == consts.PaymentFailed {
		r.Response.Writeln("Payment Failure")
	} else if one.Status < consts.PaymentSuccess {
		r.Response.Writeln("Payment Already Success")
	} else if one.ExpireTime != 0 && one.ExpireTime < gtime.Now().Timestamp() {
		r.Response.Writeln("Link Expired")
	} else if len(one.GatewayLink) > 0 {
		r.Response.RedirectTo(one.GatewayLink)
	} else if strings.Contains(one.Link, "unibee.top") {
		r.Response.Writeln("Server Error")
	} else {
		r.Response.RedirectTo(one.Link) // old version
	}
}
