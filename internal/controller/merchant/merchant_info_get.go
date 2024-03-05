package merchant

import (
	"context"
	"unibee/api/merchant/info"
	_interface "unibee/internal/interface"
	"unibee/internal/query"
)

func (c *ControllerMerchantinfo) Get(ctx context.Context, req *info.GetReq) (res *info.GetRes, err error) {
	return &info.GetRes{Merchant: query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))}, nil
}
