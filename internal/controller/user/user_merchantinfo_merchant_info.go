package user

import (
	"context"
	"unibee/api/user/merchantinfo"
	_interface "unibee/internal/interface"
	"unibee/internal/query"
)

func (c *ControllerMerchantinfo) MerchantInfo(ctx context.Context, req *merchantinfo.MerchantInfoReq) (res *merchantinfo.MerchantInfoRes, err error) {
	return &merchantinfo.MerchantInfoRes{MerchantInfo: query.GetMerchantInfoById(ctx, _interface.GetMerchantId(ctx))}, nil
}
