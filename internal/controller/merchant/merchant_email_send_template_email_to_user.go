package merchant

import (
	"context"
	"fmt"
	_interface "unibee/internal/interface/context"
	email2 "unibee/internal/logic/email"
	"unibee/internal/logic/invoice/handler"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/email"
)

func (c *ControllerEmail) SendTemplateEmailToUser(ctx context.Context, req *email.SendTemplateEmailToUserReq) (res *email.SendTemplateEmailToUserRes, err error) {
	user := query.GetUserAccountById(ctx, uint64(req.UserId))
	utility.Assert(user != nil, "User not found")
	utility.Assert(user.MerchantId == _interface.GetMerchantId(ctx), "User merchant not match")
	utility.Assert(len(req.TemplateName) > 0, "Invalid Template")

	var pdfFileName string
	if len(req.AttachInvoiceId) == 0 && req.Variables != nil && req.Variables["AttachInvoiceId"] != nil {
		req.AttachInvoiceId = fmt.Sprintf("%s", req.Variables["AttachInvoiceId"])
	}
	if len(req.AttachInvoiceId) > 0 {
		one := query.GetInvoiceByInvoiceId(ctx, req.AttachInvoiceId)
		utility.Assert(one != nil, "invoice not found")
		utility.Assert(one.UserId > 0 && int64(one.UserId) == req.UserId, "invoice userId not match")
		pdfFileName = handler.GenerateInvoicePdf(ctx, one)
	}
	err = email2.SendTemplateEmailByOpenApi(ctx, _interface.GetMerchantId(ctx), user.Email, user.TimeZone, user.Language, req.TemplateName, pdfFileName, req.Variables)
	if err != nil {
		return nil, err
	}
	return &email.SendTemplateEmailToUserRes{}, nil
}
