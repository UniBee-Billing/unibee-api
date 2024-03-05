package merchant

import (
	"context"
	"unibee/api/merchant/user"
	"unibee/internal/logic/auth"
)

func (c *ControllerUser) Release(ctx context.Context, req *user.ReleaseReq) (res *user.ReleaseRes, err error) {
	auth.ReleaseUser(ctx, req.UserId)
	return &user.ReleaseRes{}, nil
}
