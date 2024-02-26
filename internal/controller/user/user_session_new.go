package user

import (
	"context"
	session2 "unibee/internal/logic/session"
	"unibee/utility"

	"unibee/api/user/session"
)

func (c *ControllerSession) New(ctx context.Context, req *session.NewReq) (res *session.NewRes, err error) {
	utility.Assert(len(req.Email) > 0, "email is nil")
	utility.Assert(len(req.ExternalUserId) > 0, "externalUserId is nil")
	return session2.NewUserSession(ctx, req)
}
