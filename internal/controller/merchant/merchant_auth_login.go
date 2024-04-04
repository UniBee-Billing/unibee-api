package merchant

import (
	"context"
	"unibee/api/bean/detail"
	"unibee/api/merchant/auth"
	"unibee/internal/logic/member"
	"unibee/utility"
)

func (c *ControllerAuth) Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error) {
	utility.Assert(req.Email != "", "Email Cannot Be Empty")
	utility.Assert(req.Password != "", "Password Cannot Be Empty")
	one, token := member.PasswordLogin(ctx, req.Email, req.Password)
	return &auth.LoginRes{MerchantMember: detail.ConvertMemberToDetail(ctx, one), Token: token}, nil

}
