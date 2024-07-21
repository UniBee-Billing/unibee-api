package merchant

import (
	"context"
	"unibee/api/merchant/user"
	user2 "unibee/internal/logic/user"
)

func (c *ControllerUser) Frozen(ctx context.Context, req *user.FrozenReq) (res *user.FrozenRes, err error) {
	user2.FrozenUser(ctx, req.UserId)
	return &user.FrozenRes{}, nil
}
