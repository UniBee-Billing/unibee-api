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
	//isDuplicatedInvoke := false
	//defer func() {
	//	if !isDuplicatedInvoke {
	//		utility.ReleaseLock(ctx, redisKey)
	//	}
	//}()

	if !utility.TryLock(ctx, redisKey, 10) {
		//isDuplicatedInvoke = true
		utility.Assert(false, "click too fast, please wait for second")
	}

	verificationCode := generateRandomString(6)
	fmt.Printf("verification ", verificationCode)
	_, err = g.Redis().Set(ctx, req.Email+"-merchant-verify", verificationCode)
	if err != nil {
		// return nil, gerror.New("internal error")
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	_, err = g.Redis().Expire(ctx, req.Email+"-merchant-verify", 3*60)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	//email.SendEmailToUser(req.Email, "Login Code for "+req.Email+" from UniBee", verificationCode)
	merchantUser := query.GetMerchantAccountByEmail(ctx, req.Email)
	utility.Assert(merchantUser != nil, "merchant user not found")
	err = email.SendTemplateEmail(ctx, 0, req.Email, "", email.TemplateUserOTPLogin, "", &email.TemplateVariable{
		UserName:         merchantUser.FirstName + " " + merchantUser.LastName,
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	if err != nil {
		return nil, err
	}
	return &auth.LoginOtpRes{}, nil
}
