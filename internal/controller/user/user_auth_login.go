package user

import (
	"context"
	"fmt"
	"unibee-api/api/user/auth"
	"unibee-api/internal/logic/jwt"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerAuth) Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error) {
	utility.Assert(req.Email != "", "email cannot be empty")
	utility.Assert(req.Password != "", "password cannot be empty")

	var one *entity.UserAccount
	one = query.GetUserAccountByEmail(ctx, req.Email)
	utility.Assert(one != nil, "Login Failed")
	utility.Assert(one.Status == 0, "account status abnormal")
	utility.Assert(utility.ComparePasswords(one.Password, req.Password), "Login Failed, Password Not Match")

	token, err := jwt.CreatePortalToken(jwt.TOKENTYPEUSER, one.MerchantId, one.Id, req.Email)
	fmt.Println("logged-in, save email/id in token: ", req.Email, "/", one.Id)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", one.Id)), "Cache Error")
	one.Password = ""
	return &auth.LoginRes{User: one, Token: token}, nil
}
