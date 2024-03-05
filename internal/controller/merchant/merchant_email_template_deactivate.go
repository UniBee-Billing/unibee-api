package merchant

import (
	"context"
	"unibee/api/merchant/email"
	_interface "unibee/internal/interface"
	email2 "unibee/internal/logic/email"
)

func (c *ControllerEmail) TemplateDeactivate(ctx context.Context, req *email.TemplateDeactivateReq) (res *email.TemplateDeactivateRes, err error) {
	err = email2.DeactivateMerchantEmailTemplate(ctx, _interface.GetMerchantId(ctx), req.TemplateName)
	if err != nil {
		return nil, err
	} else {
		return &email.TemplateDeactivateRes{}, nil
	}
}
