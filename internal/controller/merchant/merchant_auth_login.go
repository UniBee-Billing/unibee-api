package merchant

import (
	"context"
	"fmt"
	"unibee-api/api/merchant/auth"
	_interface "unibee-api/internal/interface"
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

	var newOne *entity.MerchantUserAccount
	newOne = query.GetMerchantUserAccountByEmail(ctx, _interface.BizCtx().Get(ctx).MerchantId, req.Email)
	utility.Assert(newOne != nil, "Login Failed")
	utility.Assert(utility.ComparePasswords(newOne.Password, req.Password), "Login Failed, Password Not Match")

	token, err := jwt.CreatePortalToken(jwt.TOKENTYPEMERCHANTUSER, newOne.MerchantId, newOne.Id, req.Email)
	fmt.Println("logged-in, save email/id in token: ", req.Email, "/", newOne.Id)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("MerchantUser#%d", newOne.Id)), "Cache Error")
	newOne.Password = ""
	return &auth.LoginRes{MerchantUser: newOne, Token: token}, nil

}
