package session

import (
	"context"
	"unibee-api/internal/logic/session"
	"unibee-api/utility"

	"unibee-api/api/session/user"
)

func (c *ControllerUser) New(ctx context.Context, req *user.NewReq) (res *user.NewRes, err error) {
	utility.Assert(len(req.Email) > 0, "email is nil")
	utility.Assert(len(req.ExternalUserId) > 0, "externalUserId is nil")
	return session.NewUserSession(ctx, req)
}
