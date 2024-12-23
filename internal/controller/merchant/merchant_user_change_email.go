package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	user2 "unibee/internal/logic/user"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/user"
)

func (c *ControllerUser) ChangeEmail(ctx context.Context, req *user.ChangeEmailReq) (res *user.ChangeEmailRes, err error) {
	utility.Assert(len(req.NewEmail) > 0, "Invalid Email")
	utility.Assert(len(req.ExternalUserId) > 0 || req.UserId > 0, "either ExternalUserId or UserId needed")
	if req.UserId == 0 && len(req.ExternalUserId) > 0 {
		one := query.GetUserAccountByExternalUserId(ctx, _interface.GetMerchantId(ctx), req.ExternalUserId)
		utility.Assert(one != nil, "can't find user by ExternalUserId")
		req.UserId = one.Id
	}
	user2.ChangeUserEmail(ctx, req.UserId, req.NewEmail)
	return &user.ChangeEmailRes{}, nil
}
