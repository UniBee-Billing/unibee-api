package user

import (
	"context"
	"fmt"
	"unibee/api/user/auth"
	"unibee/internal/cmd/i18n"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/email"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerAuth) LoginOtp(ctx context.Context, req *auth.LoginOtpReq) (res *auth.LoginOtpRes, err error) {

	redisKey := fmt.Sprintf("UserAuth-Login-Email:%s", req.Email)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, i18n.LocalizationFormat(ctx, "{#ClickTooFast}"))
	}

	verificationCode := utility.GenerateRandomCode(6)
	fmt.Printf("verification %s\n", verificationCode)
	_, err = g.Redis().Set(ctx, req.Email+"-Verify", verificationCode)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, req.Email+"-verify", 3*60)
	utility.AssertError(err, "Server Error")

	user := query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
	utility.Assert(user != nil, "user not found")
	utility.Assert(user.Status != 2, "Your account has been suspended. Please contact billing admin for further assistance.")
	err = email.SendTemplateEmail(ctx, user.MerchantId, req.Email, user.TimeZone, user.Language, email.TemplateUserOTPLogin, "", &email.TemplateVariable{
		UserName:         user.FirstName + " " + user.LastName,
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	utility.AssertError(err, "Server Error")
	return &auth.LoginOtpRes{}, nil
}
