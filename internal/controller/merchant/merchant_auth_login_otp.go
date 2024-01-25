package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/email"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/merchant/auth"
)

func (c *ControllerAuth) LoginOtp(ctx context.Context, req *auth.LoginOtpReq) (res *auth.LoginOtpRes, err error) {
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
	err = email.SendTemplateEmail(ctx, 0, req.Email, email.TemplateUserOTPLogin, "", &email.TemplateVariable{
		UserName:         merchantUser.UserName,
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	if err != nil {
		return nil, err
	}
	return &auth.LoginOtpRes{}, nil
}
