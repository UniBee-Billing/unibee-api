package jwt

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"net/http"
)

func AppendRequestCookieWithToken(ctx context.Context, token string) {
	g.RequestFromCtx(ctx).AddCookie(&http.Cookie{
		Name:  "_UniBeeCookie",
		Value: token,
		//Path:     "/",
		Domain:   "localhost",
		Expires:  gtime.Now().AddDate(0, 1, 1).Time,
		MaxAge:   0,
		Secure:   false,
		HttpOnly: false,
		SameSite: 0,
		Raw:      "",
	})
}
