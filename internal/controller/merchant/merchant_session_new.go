package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	session2 "unibee/internal/logic/session"
	"unibee/utility"

	"unibee/api/merchant/session"
)

func (c *ControllerSession) New(ctx context.Context, req *session.NewReq) (res *session.NewRes, err error) {
	utility.Assert(len(req.Email) > 0, "email is nil")
	utility.Assert(len(req.ExternalUserId) > 0, "externalUserId is nil")
	return session2.NewUserSession(ctx, _interface.GetMerchantId(ctx), req)
}
