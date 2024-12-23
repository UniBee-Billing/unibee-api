package merchant

import (
	"context"
	"unibee/api/merchant/email"
	_interface "unibee/internal/interface/context"
	email2 "unibee/internal/logic/email"
)

func (c *ControllerEmail) TemplateSetDefault(ctx context.Context, req *email.TemplateSetDefaultReq) (res *email.TemplateSetDefaultRes, err error) {
	err = email2.SetMerchantEmailTemplateDefault(ctx, _interface.GetMerchantId(ctx), req.TemplateName)
	if err != nil {
		return nil, err
	} else {
		return &email.TemplateSetDefaultRes{}, nil
	}
}
