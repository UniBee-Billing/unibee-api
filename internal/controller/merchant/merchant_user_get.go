package merchant

import (
	"context"
	"go-oversea-pay/internal/query"

	"go-oversea-pay/api/merchant/user"
)

func (c *ControllerUser) Get(ctx context.Context, req *user.GetReq) (res *user.GetRes, err error) {
	return &user.GetRes{User: query.GetUserAccountById(ctx, uint64(req.UserId))}, nil
}
