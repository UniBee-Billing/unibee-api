package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"unibee-api/internal/logic/email"
	"unibee-api/internal/query"
	"unibee-api/utility"

	"unibee-api/api/merchant/auth"
)

func (c *ControllerAuth) PasswordForgetOtp(ctx context.Context, req *auth.PasswordForgetOtpReq) (res *auth.PasswordForgetOtpRes, err error) {
	redisKey := fmt.Sprintf("MerchantAuth-PasswordForgetOtp-Email:%s", req.Email)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}

	verificationCode := utility.GenerateRandomCode(6)
	fmt.Printf("verification %s", verificationCode)
	_, err = g.Redis().Set(ctx, req.Email+"-MerchantAuth-PasswordForgetOtp-Verify", verificationCode)
	if err != nil {
		// return nil, gerror.New("internal error")
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	_, err = g.Redis().Expire(ctx, req.Email+"-MerchantAuth-PasswordForgetOtp-Verify", 3*60)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	merchantUser := query.GetMerchantUserAccountByEmail(ctx, req.Email)
	utility.Assert(merchantUser != nil, "merchantUser not found")
	err = email.SendTemplateEmail(ctx, merchantUser.MerchantId, req.Email, "", email.TemplateUserOTPLogin, "", &email.TemplateVariable{
		UserName:         merchantUser.FirstName + " " + merchantUser.LastName,
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	utility.AssertError(err, "Server Error")
	return &auth.PasswordForgetOtpRes{}, nil
}