package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	email2 "unibee/internal/logic/email"
	"unibee/utility"

	"unibee/api/merchant/email"
)

func (c *ControllerEmail) TemplateActivate(ctx context.Context, req *email.TemplateActivateReq) (res *email.TemplateActivateRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}
	err = email2.ActivateMerchantEmailTemplate(ctx, _interface.GetMerchantId(ctx), req.TemplateName)
	if err != nil {
		return nil, err
	} else {
		return &email.TemplateActivateRes{}, nil
	}
}