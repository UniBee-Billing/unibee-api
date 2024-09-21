package merchant

import (
	"context"
	"unibee/api/merchant/email"
	_interface "unibee/internal/interface"
	email2 "unibee/internal/logic/email"
)

func (c *ControllerEmail) TemplateUpdate(ctx context.Context, req *email.TemplateUpdateReq) (res *email.TemplateUpdateRes, err error) {
	err = email2.UpdateMerchantEmailTemplate(ctx, _interface.GetMerchantId(ctx), req.TemplateName, req.TemplateTitle, req.TemplateContent)
	if err != nil {
		return nil, err
	} else {
		return &email.TemplateUpdateRes{}, nil
	}
}
