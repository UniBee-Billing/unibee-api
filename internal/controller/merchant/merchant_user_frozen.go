package merchant

import (
	"context"
	"unibee/api/merchant/user"
	"unibee/internal/logic/auth"
)

func (c *ControllerUser) Frozen(ctx context.Context, req *user.FrozenReq) (res *user.FrozenRes, err error) {
	auth.FrozenUser(ctx, req.UserId)
	return &user.FrozenRes{}, nil
}
