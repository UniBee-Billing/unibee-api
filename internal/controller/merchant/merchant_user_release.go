package merchant

import (
	"context"
	"unibee/api/merchant/user"
	user2 "unibee/internal/logic/user"
)

func (c *ControllerUser) Release(ctx context.Context, req *user.ReleaseReq) (res *user.ReleaseRes, err error) {
	user2.ReleaseUser(ctx, req.UserId)
	return &user.ReleaseRes{}, nil
}
