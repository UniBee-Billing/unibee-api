package _interface

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee-api/internal/model"
)

type IBizCtx interface {
	Init(r *ghttp.Request, customCtx *model.Context)
	Get(ctx context.Context) *model.Context
	SetUser(ctx context.Context, ctxUser *model.ContextUser)
	SetMerchantUser(ctx context.Context, ctxMerchantUser *model.ContextMerchantUser)
	SetData(ctx context.Context, data g.Map)
}

var localBizCtx IBizCtx

func BizCtx() IBizCtx {
	if localBizCtx == nil {
		panic("implement not found for interface IBizCtx, forgot register?")
	}
	return localBizCtx
}

const (
	SystemAssertPrefix = "system_assert: "
)

func GetMerchantId(ctx context.Context) uint64 {
	if BizCtx().Get(ctx).MerchantId <= 0 {
		panic(SystemAssertPrefix + "Invalid Merchant")
	}
	return BizCtx().Get(ctx).MerchantId
}

func RegisterBizCtx(i IBizCtx) {
	localBizCtx = i
}
