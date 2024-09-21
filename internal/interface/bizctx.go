package _interface

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee/internal/model"
)

type IContext interface {
	Init(r *ghttp.Request, customCtx *model.Context)
	Get(ctx context.Context) *model.Context
	SetUser(ctx context.Context, ctxUser *model.ContextUser)
	SetMerchantMember(ctx context.Context, ctxMerchantMember *model.ContextMerchantMember)
	SetData(ctx context.Context, data g.Map)
}

var singleTonContext IContext

func Context() IContext {
	if singleTonContext == nil {
		panic("implement not found for interface IContext, forgot register?")
	}
	return singleTonContext
}

const (
	SystemAssertPrefix = "system_assert: "
)

func GetMerchantId(ctx context.Context) uint64 {
	if Context().Get(ctx).MerchantId <= 0 {
		panic(SystemAssertPrefix + "Invalid Merchant")
	}
	return Context().Get(ctx).MerchantId
}

func RegisterContext(i IContext) {
	singleTonContext = i
}
