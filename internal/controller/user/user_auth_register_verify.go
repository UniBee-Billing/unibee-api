package user

import (
	"context"
	"encoding/json"
	"unibee/api/bean"
	"unibee/api/user/auth"
	_interface "unibee/internal/interface"
	auth2 "unibee/internal/logic/user"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerAuth) RegisterVerify(ctx context.Context, req *auth.RegisterVerifyReq) (res *auth.RegisterVerifyRes, err error) {
	verificationCode, err := g.Redis().Get(ctx, CacheKeyUserRegisterPrefix+req.Email+"-verify")
	utility.AssertError(err, "Server Error")
	utility.Assert(verificationCode != nil, "Invalid Code")
	utility.Assert((verificationCode.String()) == req.VerificationCode, "Invalid Code")

	userStr, err := g.Redis().Get(ctx, CacheKeyUserRegisterPrefix+req.Email)
	utility.AssertError(err, "Server Error")
	utility.Assert(userStr != nil, "Invalid Code")
	var newReq *auth2.NewReq
	err = json.Unmarshal([]byte(userStr.String()), &newReq)
	newReq.MerchantId = _interface.GetMerchantId(ctx)

	user, err := auth2.CreateUser(ctx, newReq)
	newOne := query.GetUserAccountById(ctx, user.Id)
	utility.Assert(newOne != nil, "Server Error")
	return &auth.RegisterVerifyRes{User: bean.SimplifyUserAccount(newOne)}, nil
}
