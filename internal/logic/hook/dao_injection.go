package hook

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	dao "unibee/internal/dao/oversea_pay"
)

func DaoHookInjection(ctx context.Context) {
	dao.Subscription.Ctx(ctx).Hook(gdb.HookHandler{
		Select: func(ctx context.Context, in *gdb.HookSelectInput) (result gdb.Result, err error) {
			return
		},
		Insert: nil,
		Update: nil,
		Delete: nil,
	})
}
