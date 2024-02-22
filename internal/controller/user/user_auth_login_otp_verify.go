package user

import (
	"context"
	"fmt"
	"unibee-api/api/user/auth"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/jwt"
	"unibee-api/utility"

	// "github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
)

func (c *ControllerAuth) LoginOtpVerify(ctx context.Context, req *auth.LoginOtpVerifyReq) (res *auth.LoginOtpVerifyRes, err error) {
	verificationCode, err := g.Redis().Get(ctx, req.Email+"-Verify")
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	utility.Assert(verificationCode != nil, "code expired")
	utility.Assert((verificationCode.String()) == req.VerificationCode, "code not match")
	var one *entity.UserAccount
	one = query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
	utility.Assert(one != nil, "Login Failed")
	token, err := jwt.CreatePortalToken(jwt.TOKENTYPEUSER, one.MerchantId, one.Id, req.Email)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", one.Id)), "Cache Error")
	one.Password = ""
	return &auth.LoginOtpVerifyRes{User: one, Token: token}, nil
}
