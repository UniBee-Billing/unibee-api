package user

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	v1 "go-oversea-pay/api/user/profile"

	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
)

func (c *ControllerProfile) Profile(ctx context.Context, req *v1.ProfileReq) (res *v1.ProfileRes, err error) {

	var newOne *entity.UserAccount
	newOne = query.GetUserAccountById(ctx, 100)

	if newOne == nil {
		// return nil, gerror.New("internal err: user not found")
		return nil, gerror.NewCode(gcode.New(400, "login failed", nil))
	}

	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
