package user

import (
	"context"
	"unibee/api/user/merchant"
	_interface "unibee/internal/interface"
	"unibee/internal/query"
)

func (c *ControllerMerchant) Get(ctx context.Context, req *merchant.GetReq) (res *merchant.GetRes, err error) {
	return &merchant.GetRes{Merchant: query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))}, nil
}
