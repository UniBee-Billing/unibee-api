package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee-api/api/merchant/auth"
)

func (c *ControllerAuth) PasswordForgetOtp(ctx context.Context, req *auth.PasswordForgetOtpReq) (res *auth.PasswordForgetOtpRes, err error) {
	return nil, gerror.New("not support")
	//redisKey := fmt.Sprintf("MerchantAuth-PasswordForgetOtp-Email:%s", req.Email)
	//if !utility.TryLock(ctx, redisKey, 10) {
	//	utility.Assert(false, "click too fast, please wait for second")
	//}
	//
	//verificationCode := utility.GenerateRandomCode(6)
	//fmt.Printf("verification %s", verificationCode)
	//_, err = g.Redis().Set(ctx, req.Email+"-MerchantAuth-PasswordForgetOtp-Verify", verificationCode)
	//if err != nil {
	//	// return nil, gerror.New("internal error")
	//	return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	//}
	//_, err = g.Redis().Expire(ctx, req.Email+"-MerchantAuth-PasswordForgetOtp-Verify", 3*60)
	//if err != nil {
	//	return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	//}
	//
	//user := query.GetUserAccountByEmail(ctx, _interface.BizCtx().Get(ctx).MerchantId, req.Email)
	//utility.Assert(user != nil, "user not found")
	//err = email.SendTemplateEmail(ctx, 0, req.Email, "", email.TemplateUserOTPLogin, "", &email.TemplateVariable{
	//	UserName:         user.FirstName + " " + user.LastName,
	//	CodeExpireMinute: "3",
	//	Code:             verificationCode,
	//})
	//if err != nil {
	//	return nil, err
	//}
	//return &auth.PasswordForgetOtpRes{}, nil
}
