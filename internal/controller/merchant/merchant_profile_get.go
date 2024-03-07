package merchant

import (
	"context"
	"unibee/api/merchant/profile"
	_interface "unibee/internal/interface"
	"unibee/internal/query"
)

func (c *ControllerMerchantProfile) Get(ctx context.Context, req *profile.GetReq) (res *profile.GetRes, err error) {
	return &profile.GetRes{Merchant: query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))}, nil
}
