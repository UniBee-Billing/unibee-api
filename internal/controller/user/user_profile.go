package user

import (
	"context"
	"unibee/api/bean"
	v1 "unibee/api/user/profile"

	_interface "unibee/internal/interface"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/internal/query"
)

func (c *ControllerProfile) Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error) {
	var newOne = query.GetUserAccountById(ctx, _interface.BizCtx().Get(ctx).User.Id)
	if newOne == nil {
		// return nil, gerror.New("internal err: user not found")
		return nil, gerror.NewCode(gcode.New(400, "login failed", nil))
	}

	return &v1.GetRes{User: bean.SimplifyUserAccount(newOne)}, nil
}
