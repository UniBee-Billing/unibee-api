package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/auth"
	"unibee/internal/logic/merchant"
	"unibee/utility"
)

func (c *ControllerAuth) Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error) {
	utility.Assert(req.Email != "", "Email Cannot Be Empty")
	utility.Assert(req.Password != "", "Password Cannot Be Empty")
	one, token := merchant.PasswordLogin(ctx, req.Email, req.Password)
	return &auth.LoginRes{MerchantMember: bean.SimplifyMerchantMember(one), Token: token}, nil

}
