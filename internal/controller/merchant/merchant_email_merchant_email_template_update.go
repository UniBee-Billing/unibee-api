package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	email2 "unibee/internal/logic/email"
	"unibee/utility"

	"unibee/api/merchant/email"
)

func (c *ControllerEmail) MerchantEmailTemplateUpdate(ctx context.Context, req *email.MerchantEmailTemplateUpdateReq) (res *email.MerchantEmailTemplateUpdateRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}
	err = email2.UpdateMerchantEmailTemplate(ctx, _interface.GetMerchantId(ctx), req.TemplateName, req.TemplateTitle, req.TemplateContent)
	if err != nil {
		return nil, err
	} else {
		return &email.MerchantEmailTemplateUpdateRes{}, nil
	}
}
