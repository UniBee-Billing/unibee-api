package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/info"
)

func (c *ControllerMerchantinfo) Get(ctx context.Context, req *info.GetReq) (res *info.GetRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.MerchantId > 0, "MerchantId invalid")
	}
	return &info.GetRes{Merchant: query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))}, nil
}
