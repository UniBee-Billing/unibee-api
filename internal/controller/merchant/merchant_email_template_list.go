package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	email2 "unibee/internal/logic/email"
	"unibee/utility"

	"unibee/api/merchant/email"
)

func (c *ControllerEmail) TemplateList(ctx context.Context, req *email.TemplateListReq) (res *email.TemplateListRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}
	list := email2.GetMerchantEmailTemplateList(ctx, _interface.GetMerchantId(ctx))
	return &email.TemplateListRes{EmailTemplateList: list}, nil
}