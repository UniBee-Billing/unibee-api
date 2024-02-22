package user

import (
	"context"
	"fmt"
	"unibee-api/api/user/auth"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/email"
	"unibee-api/internal/query"
	"unibee-api/utility"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerAuth) LoginOtp(ctx context.Context, req *auth.LoginOtpReq) (res *auth.LoginOtpRes, err error) {

	redisKey := fmt.Sprintf("UserAuth-Login-Email:%s", req.Email)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}

	verificationCode := utility.GenerateRandomCode(6)
	fmt.Printf("verification %s", verificationCode)
	_, err = g.Redis().Set(ctx, req.Email+"-Verify", verificationCode)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, req.Email+"-verify", 3*60)
	utility.AssertError(err, "Server Error")

	user := query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
	utility.Assert(user != nil, "user not found")
	utility.Assert(user.Status == 0, "account status abnormal")
	err = email.SendTemplateEmail(ctx, user.MerchantId, req.Email, "", email.TemplateUserOTPLogin, "", &email.TemplateVariable{
		UserName:         user.FirstName + " " + user.LastName,
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	utility.AssertError(err, "Server Error")
	return &auth.LoginOtpRes{}, nil
}
