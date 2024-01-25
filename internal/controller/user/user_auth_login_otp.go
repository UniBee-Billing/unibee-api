package user

import (
	"context"
	"fmt"
	"go-oversea-pay/api/user/auth"
	"go-oversea-pay/internal/logic/email"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	// "github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerAuth) LoginOtp(ctx context.Context, req *auth.LoginOtpReq) (res *auth.LoginOtpRes, err error) {
	verificationCode := generateRandomString(6)
	fmt.Printf("verification ", verificationCode)
	_, err = g.Redis().Set(ctx, req.Email+"-verify", verificationCode)
	if err != nil {
		// return nil, gerror.New("internal error")
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	_, err = g.Redis().Expire(ctx, req.Email+"-verify", 3*60)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	//email.SendEmailToUser(req.Email, "Login Code for "+req.Email+" from UniBee", verificationCode)
	user := query.GetUserAccountByEmail(ctx, req.Email)
	utility.Assert(user != nil, "user not found")
	err = email.SendTemplateEmail(ctx, 0, req.Email, email.TemplateUserOTPLogin, "", &email.TemplateVariable{
		UserName:         user.UserName,
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	if err != nil {
		return nil, err
	}
	return &auth.LoginOtpRes{}, nil
}
