package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	email2 "unibee-api/internal/logic/email"
	"unibee-api/utility"

	"unibee-api/api/merchant/email"
)

func (c *ControllerEmail) MerchantEmailTemplateSetDefault(ctx context.Context, req *email.MerchantEmailTemplateSetDefaultReq) (res *email.MerchantEmailTemplateSetDefaultRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	err = email2.SetMerchantEmailTemplateDefault(ctx, _interface.GetMerchantId(ctx), req.TemplateName)
	if err != nil {
		return nil, err
	} else {
		return &email.MerchantEmailTemplateSetDefaultRes{}, nil
	}
}
