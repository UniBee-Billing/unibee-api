package user

import (
	"context"
	"fmt"
	"unibee-api/api/user/auth"
	auth2 "unibee-api/internal/logic/auth"
	"unibee-api/utility"

	// "github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
)

func (c *ControllerAuth) LoginOtpVerify(ctx context.Context, req *auth.LoginOtpVerifyReq) (res *auth.LoginOtpVerifyRes, err error) {
	verificationCode, err := g.Redis().Get(ctx, req.Email+"-verify")
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	utility.Assert(verificationCode != nil, "code expired")
	//if verificationCode == nil { // expired
	//	return nil, gerror.NewCode(gcode.New(400, "code expired", nil))
	//}
	utility.Assert((verificationCode.String()) == req.VerificationCode, "code not match")
	//if (verificationCode.String()) != req.VerificationCode {
	//	return nil, gerror.NewCode(gcode.New(400, "code not match", nil))
	//}

	var newOne *entity.UserAccount
	newOne = query.GetUserAccountByEmail(ctx, req.Email)
	utility.Assert(newOne != nil, "Login Failed")
	//if newOne == nil {
	//	return nil, gerror.NewCode(gcode.New(400, "login failed", nil))
	//}

	token, err := auth2.CreateToken(req.Email, newOne.Id)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	utility.Assert(auth2.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", newOne.Id)), "Cache Error")
	newOne.Password = ""
	return &auth.LoginOtpVerifyRes{User: newOne, Token: token}, nil
}
