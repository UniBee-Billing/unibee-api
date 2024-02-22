package user

import (
	"context"
	"fmt"
	auth2 "unibee-api/internal/logic/auth"
	"unibee-api/internal/logic/jwt"
	"unibee-api/utility"

	"unibee-api/api/user/auth"
)

func (c *ControllerAuth) SessionLogin(ctx context.Context, req *auth.SessionLoginReq) (res *auth.SessionLoginRes, err error) {
	one := auth2.UserSessionTransfer(ctx, req.Session)
	utility.Assert(one != nil, "Login Failed")

	token, err := jwt.CreatePortalToken(jwt.TOKENTYPEUSER, one.MerchantId, one.Id, one.Email)
	fmt.Println("session logged-in, save email/id in token: ", one.Email, "/", one.Id)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", one.Id)), "Cache Error")
	one.Password = ""
	return &auth.SessionLoginRes{User: one, Token: token}, nil
}
