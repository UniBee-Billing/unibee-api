package merchant

import (
	"context"
	"unibee/api/bean/detail"
	"unibee/api/merchant/user"
	_interface "unibee/internal/interface"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerUser) Get(ctx context.Context, req *user.GetReq) (res *user.GetRes, err error) {
	one := query.GetUserAccountById(ctx, uint64(req.UserId))
	utility.Assert(one.MerchantId == _interface.GetMerchantId(ctx), "wrong merchant account")
	return &user.GetRes{User: detail.ConvertUserAccountToDetail(ctx, one)}, nil
}
