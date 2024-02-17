package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee-api/internal/controller/session"
)

func UserSession(ctx context.Context, group *ghttp.RouterGroup) {
	// auth
	group.Group("/user", func(group *ghttp.RouterGroup) {
		group.Bind(
			session.NewUser(),
		)
	})
}
