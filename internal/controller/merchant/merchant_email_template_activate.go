package merchant

import (
	"context"
	"unibee/api/merchant/email"
	_interface "unibee/internal/interface/context"
	email2 "unibee/internal/logic/email"
)

func (c *ControllerEmail) TemplateActivate(ctx context.Context, req *email.TemplateActivateReq) (res *email.TemplateActivateRes, err error) {
	err = email2.ActivateMerchantEmailTemplate(ctx, _interface.GetMerchantId(ctx), req.TemplateName)
	if err != nil {
		return nil, err
	} else {
		return &email.TemplateActivateRes{}, nil
	}
}
