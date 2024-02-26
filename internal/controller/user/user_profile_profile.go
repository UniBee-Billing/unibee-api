package user

import (
	"context"
	v1 "unibee/api/user/profile"

	_interface "unibee/internal/interface"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

func (c *ControllerProfile) Profile(ctx context.Context, req *v1.ProfileReq) (res *v1.ProfileRes, err error) {
	var newOne *entity.UserAccount = query.GetUserAccountById(ctx, _interface.BizCtx().Get(ctx).User.Id)
	if newOne == nil {
		// return nil, gerror.New("internal err: user not found")
		return nil, gerror.NewCode(gcode.New(400, "login failed", nil))
	}
	
	return &v1.ProfileRes{User: newOne}, nil
}
 