package user

import (
	"context"
	"fmt"
	"unibee-api/api/user/auth"
	"unibee-api/internal/logic/email"
	"unibee-api/internal/query"
	"unibee-api/utility"

	// "github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerAuth) LoginOtp(ctx context.Context, req *auth.LoginOtpReq) (res *auth.LoginOtpRes, err error) {

	redisKey := fmt.Sprintf("UserAuth-Login-Email:%s", req.Email)
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
		UserName:         user.FirstName + " " + user.LastName,
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	if err != nil {
		return nil, err
	}
	return &auth.LoginOtpRes{}, nil
}
