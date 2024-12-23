package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	email2 "unibee/internal/logic/email"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/email"
)

func (c *ControllerEmail) SendTemplateEmailToUser(ctx context.Context, req *email.SendTemplateEmailToUserReq) (res *email.SendTemplateEmailToUserRes, err error) {
	user := query.GetUserAccountById(ctx, uint64(req.UserId))
	utility.Assert(user != nil, "User not found")
	utility.Assert(user.MerchantId == _interface.GetMerchantId(ctx), "User merchant not match")
	utility.Assert(len(req.TemplateName) > 0, "Invalid Template")

	err = email2.SendTemplateEmailByOpenApi(ctx, _interface.GetMerchantId(ctx), user.Email, user.TimeZone, user.Language, req.TemplateName, "", req.Variables)
	if err != nil {
		return nil, err
	}
	return &email.SendTemplateEmailToUserRes{}, nil
}
