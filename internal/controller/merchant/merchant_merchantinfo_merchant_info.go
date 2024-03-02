package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/merchantinfo"
)

func (c *ControllerMerchantinfo) MerchantInfo(ctx context.Context, req *merchantinfo.MerchantInfoReq) (res *merchantinfo.MerchantInfoRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.MerchantId > 0, "MerchantId invalid")
	}
	return &merchantinfo.MerchantInfoRes{MerchantInfo: query.GetMerchantById(ctx, _interface.BizCtx().Get(ctx).MerchantMember.MerchantId)}, nil
}
