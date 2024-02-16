package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/query"
	"unibee-api/utility"

	"unibee-api/api/merchant/merchantinfo"
)

func (c *ControllerMerchantinfo) MerchantInfo(ctx context.Context, req *merchantinfo.MerchantInfoReq) (res *merchantinfo.MerchantInfoRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId > 0, "MerchantId invalid")
	}
	return &merchantinfo.MerchantInfoRes{MerchantInfo: query.GetMerchantInfoById(ctx, int64(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId))}, nil
}
