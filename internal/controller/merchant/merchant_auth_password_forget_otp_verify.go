package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee-api/api/merchant/auth"
)

func (c *ControllerAuth) PasswordForgetOtpVerify(ctx context.Context, req *auth.PasswordForgetOtpVerifyReq) (res *auth.PasswordForgetOtpVerifyRes, err error) {
	//verificationCode, err := g.Redis().Get(ctx, req.Email+"-MerchantAuth-PasswordForgetOtp-Verify")
	//if err != nil {
	//	return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	//}
	//utility.Assert(verificationCode != nil, "code expired")
	//utility.Assert((verificationCode.String()) == req.VerificationCode, "code not match")
	//
	//var newOne *entity.MerchantUserAccount
	//newOne = query.GetMerchantAccountByEmail(ctx, _interface.BizCtx().Get(ctx).MerchantId, req.Email)
	//utility.Assert(newOne != nil, "User Not Found")

	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
