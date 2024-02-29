package merchant

import (
	"context"
	"fmt"
	"unibee/api/merchant/auth"
	"unibee/internal/logic/jwt"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerAuth) Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error) {
	utility.Assert(req.Email != "", "Email Cannot Be Empty")
	utility.Assert(req.Password != "", "Password Cannot Be Empty")

	var newOne *entity.MerchantUserAccount
	newOne = query.GetMerchantUserAccountByEmail(ctx, req.Email)
	utility.Assert(newOne != nil, "Email Not Found")
	utility.Assert(utility.ComparePasswords(newOne.Password, req.Password), "Login Failed, Password Not Match")

	token, err := jwt.CreatePortalToken(jwt.TOKENTYPEMERCHANTUSER, newOne.MerchantId, newOne.Id, req.Email)
	fmt.Println("logged-in, save email/id in token: ", req.Email, "/", newOne.Id)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "Server Error", nil))
	}
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("MerchantUser#%d", newOne.Id)), "Cache Error")
	newOne.Password = ""
	return &auth.LoginRes{MerchantUser: newOne, Token: token}, nil

}
