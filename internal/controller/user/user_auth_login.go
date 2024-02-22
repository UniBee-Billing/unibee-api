package user

import (
	"context"
	"fmt"
	"unibee-api/api/user/auth"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/jwt"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func (c *ControllerAuth) Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error) {
	utility.Assert(req.Email != "", "email cannot be empty")
	utility.Assert(req.Password != "", "password cannot be empty")

	var one *entity.UserAccount
	one = query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
	utility.Assert(one != nil, "Login Failed")
	utility.Assert(one.Status == 0, "account status abnormal")
	utility.Assert(utility.ComparePasswords(one.Password, req.Password), "Login Failed, Password Not Match")

	token, err := jwt.CreatePortalToken(jwt.TOKENTYPEUSER, one.MerchantId, one.Id, req.Email)
	fmt.Println("logged-in, save email/id in token: ", req.Email, "/", one.Id)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", one.Id)), "Cache Error")
	one.Password = ""
	return &auth.LoginRes{User: one, Token: token}, nil
}
