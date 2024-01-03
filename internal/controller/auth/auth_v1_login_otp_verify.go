package auth

import (
	"context"

	// "github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "go-oversea-pay/api/auth/v1"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
)

func (c *ControllerV1) LoginOtpVerify(ctx context.Context, req *v1.LoginOtpVerifyReq) (res *v1.LoginOtpVerifyRes, err error) {
	verificationCode, err := g.Redis().Get(ctx, req.Email + "-verify")
		if err != nil {
			return nil, gerror.NewCode(gcode.New(500, "server error", nil))
		}
		if verificationCode == nil { // expired
			return nil, gerror.NewCode(gcode.New(400, "code expired", nil))
		}

		if (verificationCode.String()) != req.VerificationCode {
			return nil, gerror.NewCode(gcode.New(400, "code not match", nil))
		}
	
		var newOne *entity.UserAccount
		newOne = query.GetUserAccountByEmail(ctx, req.Email)
		if newOne == nil {
			// return nil, gerror.New("internal err: user not found")
			return nil, gerror.NewCode(gcode.New(400, "login failed", nil))
		}
	
		token, err := createToken(req.Email)
		if err != nil {
			return nil, gerror.NewCode(gcode.New(500, "server error", nil))
		}

	return &v1.LoginOtpVerifyRes{User: newOne, Token: token}, nil
}
