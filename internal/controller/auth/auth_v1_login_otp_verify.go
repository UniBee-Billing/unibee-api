package auth

import (
	"context"

	// "github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "go-oversea-pay/api/auth/v1"
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

	return &v1.LoginOtpVerifyRes{}, nil
}
