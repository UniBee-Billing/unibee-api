package merchant

import (
	"context"
	"unibee-api/internal/query"

	"unibee-api/api/merchant/user"
)

func (c *ControllerUser) Get(ctx context.Context, req *user.GetReq) (res *user.GetRes, err error) {
	return &user.GetRes{User: query.GetUserAccountById(ctx, uint64(req.UserId))}, nil
}
