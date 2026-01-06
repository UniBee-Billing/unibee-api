package analytics

import (
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee/internal/cmd/config"
)

func GoAnalyticsPortal(r *ghttp.Request) {
	r.Response.RedirectTo(fmt.Sprintf("%s/analytics", config.GetConfigInstance().Server.GetServerPath()))
}
