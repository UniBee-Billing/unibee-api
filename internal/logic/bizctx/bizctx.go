package bizctx

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type sBizCtx struct{}

func init() {
	_interface.RegisterBizCtx(New())
}

func New() *sBizCtx {
	return &sBizCtx{}
}

func (s *sBizCtx) Init(r *ghttp.Request, customCtx *model.Context) {
	r.SetCtxVar(consts.ContextKey, customCtx)
}

func (s *sBizCtx) Get(ctx context.Context) *model.Context {
	value := ctx.Value(consts.ContextKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*model.Context); ok {
		return localCtx
	}
	return nil
}

func (s *sBizCtx) SetUser(ctx context.Context, ctxUser *model.ContextUser) {
	s.Get(ctx).User = ctxUser
}

func (s *sBizCtx) SetMerchantUser(ctx context.Context, ctxMerchantUser *model.ContextMerchantUser) {
	s.Get(ctx).MerchantUser = ctxMerchantUser
}

func (s *sBizCtx) SetData(ctx context.Context, data g.Map) {
	s.Get(ctx).Data = data
}
