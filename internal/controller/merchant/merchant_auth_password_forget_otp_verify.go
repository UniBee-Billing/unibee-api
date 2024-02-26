package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	auth2 "unibee-api/internal/logic/auth"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"

	"unibee-api/api/merchant/auth"
)

func (c *ControllerAuth) PasswordForgetOtpVerify(ctx context.Context, req *auth.PasswordForgetOtpVerifyReq) (res *auth.PasswordForgetOtpVerifyRes, err error) {
	verificationCode, err := g.Redis().Get(ctx, req.Email+"-MerchantAuth-PasswordForgetOtp-Verify")
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	utility.Assert(verificationCode != nil, "code expired")
	utility.Assert((verificationCode.String()) == req.VerificationCode, "code not match")

	var newOne *entity.MerchantUserAccount
	newOne = query.GetMerchantUserAccountByEmail(ctx, req.Email)
	utility.Assert(newOne != nil, "User Not Found")
	auth2.ChangeMerchantUserPasswordWithOutOldVerify(ctx, req.Email, req.NewPassword)
	return &auth.PasswordForgetOtpVerifyRes{}, nil
}
