package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	email2 "unibee-api/internal/logic/email"
	"unibee-api/utility"

	"unibee-api/api/merchant/email"
)

func (c *ControllerEmail) MerchantEmailTemplateList(ctx context.Context, req *email.MerchantEmailTemplateListReq) (res *email.MerchantEmailTemplateListRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId > 0, "merchantUserId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId == uint64(req.MerchantId), "merchantId not match")
	}
	list := email2.GetMerchantEmailTemplateList(ctx, req.MerchantId)
	return &email.MerchantEmailTemplateListRes{EmailTemplateList: list}, nil
}
