package user

import (
	"context"
	"go-oversea-pay/api/user/auth"
	// entity "go-oversea-pay/internal/model/entity/oversea_pay"
	// "go-oversea-pay/internal/query"
	// "github.com/gogf/gf/v2/errors/gcode"
	// "github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerAuth) Logout(ctx context.Context, req *auth.LogoutReq) (res *auth.LogoutRes, err error) {
	// reset ctx customer User obj
	return &auth.LogoutRes{}, nil
}
