package context

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type Context struct{}

func init() {
	_interface.RegisterContext(New())
}

func New() *Context {
	return &Context{}
}

func (s *Context) Init(r *ghttp.Request, customCtx *model.Context) {
	r.SetCtxVar(consts.ContextKey, customCtx)
}

func (s *Context) Get(ctx context.Context) *model.Context {
	value := ctx.Value(consts.ContextKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*model.Context); ok {
		return localCtx
	}
	return nil
}

func (s *Context) SetUser(ctx context.Context, ctxUser *model.ContextUser) {
	s.Get(ctx).User = ctxUser
}

func (s *Context) SetMerchantMember(ctx context.Context, ctxMerchantMember *model.ContextMerchantMember) {
	s.Get(ctx).MerchantMember = ctxMerchantMember
}

func (s *Context) SetData(ctx context.Context, data g.Map) {
	s.Get(ctx).Data = data
}
