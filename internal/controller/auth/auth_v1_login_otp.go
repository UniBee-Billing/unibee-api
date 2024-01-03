package auth

import (
	"context"
	"fmt"

	// "github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "go-oversea-pay/api/auth/v1"
	"go-oversea-pay/internal/logic/email"
)

func (c *ControllerV1) LoginOtp(ctx context.Context, req *v1.LoginOtpReq) (res *v1.LoginOtpRes, err error) {
	verificationCode := generateRandomString(6)
	fmt.Printf("verification ", verificationCode)
	_, err = g.Redis().Set(ctx, req.Email + "-verify", verificationCode)
	if err != nil {
		// return nil, gerror.New("internal error")
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	_, err = g.Redis().Expire(ctx, req.Email + "-verify", 3*60)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	
	email.SendEmailToUser(req.Email, "Login Code for " + req.Email + " from Unibee", verificationCode)
	return &v1.LoginOtpRes{}, nil 
	// return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
