package integrations

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
	"strconv"
	"unibee/internal/logic/analysis/quickbooks"
)

func QuickBooksAuthorizationEntry(r *ghttp.Request) {
	code := r.GetQuery("code")
	realmId := r.GetQuery("realmId")
	state := r.GetQuery("state")
	if code == nil || realmId == nil || state == nil {
		g.Log().Errorf(r.Context(), "QuickBooksAuthorizationEntry no code or realmId or state")
		r.Response.Status = http.StatusForbidden
		r.Exit()
		return
	}
	if len(code.String()) == 0 || len(realmId.String()) == 0 || len(state.String()) == 0 {
		g.Log().Errorf(r.Context(), "QuickBooksAuthorizationEntry no code or realmId or state")
		r.Response.Status = http.StatusForbidden
		r.Exit()
		return
	}
	g.Log().Infof(r.Context(), "QuickBooksAuthorizationEntry code:%s, realmId:%s, state:%s", code, realmId, state)
	merchantId, err := strconv.ParseInt(state.String(), 10, 64)
	if err != nil {
		g.Log().Errorf(r.Context(), "QuickBooksAuthorizationEntry invalid merchantId")
		r.Response.Status = http.StatusForbidden
		r.Exit()
		return
	}
	var returnUrl = ""
	config := quickbooks.GetMerchantQuickBooksConfig(r.Context(), uint64(merchantId))
	if config != nil && config.SetupReturnUrl != "" {
		returnUrl = config.SetupReturnUrl
	}

	err = quickbooks.SetupMerchantQuickBooksConfig(r.Context(), uint64(merchantId), &quickbooks.MerchantQuickBooksConfig{
		Code:    code.String(),
		RealmId: realmId.String(),
	})
	if err != nil {
		g.Log().Errorf(r.Context(), "SetupMerchantQuickBooksConfig err:%s", err.Error())
		r.Response.Status = http.StatusForbidden
		r.Exit()
		return
	}
	r.Response.Status = http.StatusOK
	if len(returnUrl) > 0 {
		r.Response.RedirectTo(returnUrl)
	} else {
		r.Response.Write("success")
	}
	return
}
