package merchant

import (
	"context"
	"fmt"
	"strings"
	"time"
	_interface "unibee/internal/interface/context"
	email2 "unibee/internal/logic/email"
	"unibee/internal/logic/invoice/handler"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/email"
)

func (c *ControllerEmail) SendEmailToUser(ctx context.Context, req *email.SendEmailToUserReq) (res *email.SendEmailToUserRes, err error) {
	user := query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
	utility.Assert(user != nil, "User not found")
	mailTo := strings.ToLower(req.Email)
	_, emailGatewayKey := email2.GetDefaultMerchantEmailConfigWithClusterCloud(ctx, _interface.GetMerchantId(ctx))
	if len(emailGatewayKey) == 0 {
		utility.Assert(false, "Default Email Gateway Need Setup")
	}
	var pdfFileName string
	var attachName string
	if len(req.AttachInvoiceId) == 0 && req.Variables != nil && req.Variables["AttachInvoiceId"] != nil {
		req.AttachInvoiceId = fmt.Sprintf("%s", req.Variables["AttachInvoiceId"])
	}
	if len(req.AttachInvoiceId) > 0 {
		one := query.GetInvoiceByInvoiceId(ctx, req.AttachInvoiceId)
		utility.Assert(one != nil, "invoice not found")
		utility.Assert(one.UserId > 0 && one.UserId == user.Id, "invoice userId not match")
		pdfFileName = handler.GenerateInvoicePdf(ctx, one)
		attachName = fmt.Sprintf("invoice_%s", time.Now().Format("20060102"))
	}
	err = email2.Send(ctx, &email2.SendgridEmailReq{
		MerchantId:        _interface.GetMerchantId(ctx),
		MailTo:            mailTo,
		Subject:           req.Subject,
		Content:           req.Content,
		LocalFilePath:     pdfFileName,
		AttachName:        attachName + ".pdf",
		APIKey:            emailGatewayKey,
		VariableMap:       req.Variables,
		Language:          user.Language,
		GatewayTemplateId: req.GatewayTemplateId,
	})
	if err != nil {
		return nil, err
	}
	return &email.SendEmailToUserRes{}, nil
}
