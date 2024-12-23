package user

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	_interface "unibee/internal/interface/context"
	auth2 "unibee/internal/logic/user"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/user/auth"
)

func (c *ControllerAuth) PasswordForgetOtpVerify(ctx context.Context, req *auth.PasswordForgetOtpVerifyReq) (res *auth.PasswordForgetOtpVerifyRes, err error) {
	verificationCode, err := g.Redis().Get(ctx, req.Email+"-PasswordForgetOtp-Verify")
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	utility.Assert(verificationCode != nil, "code expired")
	utility.Assert((verificationCode.String()) == req.VerificationCode, "code not match")

	auth2.ChangeUserPasswordWithOutOldVerify(ctx, _interface.GetMerchantId(ctx), req.Email, req.NewPassword)

	return &auth.PasswordForgetOtpVerifyRes{}, nil
}
