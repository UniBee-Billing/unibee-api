package merchant

import (
	"context"
	"unibee/api/merchant/email"
	_interface "unibee/internal/interface/context"
	email2 "unibee/internal/logic/email"
)

func (c *ControllerEmail) TemplateList(ctx context.Context, req *email.TemplateListReq) (res *email.TemplateListRes, err error) {
	list, total := email2.GetMerchantEmailTemplateList(ctx, _interface.GetMerchantId(ctx))
	return &email.TemplateListRes{EmailTemplateList: list, Total: total}, nil
}
