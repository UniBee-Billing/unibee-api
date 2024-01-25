package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	auth2 "go-oversea-pay/internal/logic/auth"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/merchant/auth"
)

func (c *ControllerAuth) LoginOtpVerify(ctx context.Context, req *auth.LoginOtpVerifyReq) (res *auth.LoginOtpVerifyRes, err error) {
	verificationCode, err := g.Redis().Get(ctx, req.Email+"-merchant-verify")
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	utility.Assert(verificationCode != nil, "code expired")
	//if verificationCode == nil { // expired
	//	return nil, gerror.NewCode(gcode.New(400, "code expired", nil))
	//}
	utility.Assert((verificationCode.String()) == req.VerificationCode, "code not match")
	//if (verificationCode.String()) != req.VerificationCode {
	//	return nil, gerror.NewCode(gcode.New(400, "code not match", nil))
	//}

	var newOne *entity.MerchantUserAccount
	newOne = query.GetMerchantAccountByEmail(ctx, req.Email)
	utility.Assert(newOne != nil, "Login Failed")
	//if newOne == nil {
	//	return nil, gerror.NewCode(gcode.New(400, "login failed", nil))
	//}

	token, err := createToken(req.Email, newOne.Id)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	utility.Assert(auth2.PutAuthTokenToCache(ctx, token, fmt.Sprintf("MerchantUser#%d", newOne.Id)), "Cache Error")

	return &auth.LoginOtpVerifyRes{MerchantUser: newOne, Token: token}, nil
}
