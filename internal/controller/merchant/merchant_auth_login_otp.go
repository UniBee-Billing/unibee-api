package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/email"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/auth"
)

func (c *ControllerAuth) LoginOtp(ctx context.Context, req *auth.LoginOtpReq) (res *auth.LoginOtpRes, err error) {
	redisKey := fmt.Sprintf("MerchantAuth-Login-Email:%s", req.Email)

	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}

	verificationCode := utility.GenerateRandomCode(6)
	fmt.Printf("verification :%s", verificationCode)
	_, err = g.Redis().Set(ctx, req.Email+"-MerchantAuth-Verify", verificationCode)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, req.Email+"-MerchantAuth-Verify", 3*60)
	utility.AssertError(err, "Server Error")

	merchantMember := query.GetMerchantMemberByEmail(ctx, req.Email)
	utility.Assert(merchantMember != nil, "merchant member not found")
	// deploy version todo mark send to unibee api
	err = email.SendTemplateEmail(ctx, merchantMember.MerchantId, req.Email, "", email.TemplateMerchantOTPLogin, "", &email.TemplateVariable{
		UserName:         merchantMember.FirstName + " " + merchantMember.LastName,
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	utility.AssertError(err, "Server Error")
	return &auth.LoginOtpRes{}, nil
}
