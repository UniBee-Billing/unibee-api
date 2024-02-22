package user

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/email"
	"unibee-api/internal/query"
	"unibee-api/utility"

	"unibee-api/api/user/auth"
)

func (c *ControllerAuth) PasswordForgetOtp(ctx context.Context, req *auth.PasswordForgetOtpReq) (res *auth.PasswordForgetOtpRes, err error) {
	redisKey := fmt.Sprintf("UserAuth-PasswordForgetOtp-Email:%s", req.Email)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}

	verificationCode := utility.GenerateRandomCode(6)
	fmt.Printf("verification %s", verificationCode)
	_, err = g.Redis().Set(ctx, req.Email+"-PasswordForgetOtp-Verify", verificationCode)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, req.Email+"-PasswordForgetOtp-Verify", 3*60)
	utility.AssertError(err, "Server Error")

	user := query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
	utility.Assert(user != nil, "user not found")
	utility.Assert(user.Status == 0, "account status abnormal")
	err = email.SendTemplateEmail(ctx, user.MerchantId, req.Email, "", email.TemplateUserOTPLogin, "", &email.TemplateVariable{
		UserName:         user.FirstName + " " + user.LastName,
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	if err != nil {
		return nil, err
	}
	return &auth.PasswordForgetOtpRes{}, nil
}
