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

func (c *ControllerAuth) LoginOtp(ctx context.Context, req *auth.LoginOtpReq) (res *auth.LoginOtpRes, err error) {
	redisKey := fmt.Sprintf("MerchantAuth-Login-Email:%s", req.Email)

	if !utility.TryLock(ctx, redisKey, 10) {
		//isDuplicatedInvoke = true
		utility.Assert(false, "click too fast, please wait for second")
	}

	verificationCode := utility.GenerateRandomCode(6)
	fmt.Printf("verification :%s", verificationCode)
	_, err = g.Redis().Set(ctx, req.Email+"-MerchantAuth-Verify", verificationCode)
	if err != nil {
		// return nil, gerror.New("internal error")
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	_, err = g.Redis().Expire(ctx, req.Email+"-MerchantAuth-Verify", 3*60)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	merchantUser := query.GetMerchantUserAccountByEmail(ctx, req.Email)
	utility.Assert(merchantUser != nil, "merchant user not found")
	err = email.SendTemplateEmail(ctx, merchantUser.MerchantId, req.Email, "", email.TemplateUserOTPLogin, "", &email.TemplateVariable{
		UserName:         merchantUser.FirstName + " " + merchantUser.LastName,
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	if err != nil {
		return nil, err
	}
	return &auth.LoginOtpRes{}, nil
}
