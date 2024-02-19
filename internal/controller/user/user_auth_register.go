package user

import (
	"context"
	"encoding/json"
	"fmt"
	"unibee-api/api/user/auth"
	"unibee-api/internal/logic/email"
	"unibee-api/utility"

	"github.com/gogf/gf/v2/frame/g"

	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerAuth) Register(ctx context.Context, req *auth.RegisterReq) (res *auth.RegisterRes, err error) {
	utility.Assert(len(req.Email) > 0, "Email Needed")
	utility.Assert(utility.IsEmailValid(req.Email), "Invalid Email")

	redisKey := fmt.Sprintf("UserAuth-Regist-Email:%s", req.Email)

	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}

	var newOne *entity.UserAccount
	newOne = query.GetUserAccountByEmail(ctx, req.Email) //Id(ctx, user.Id)
	utility.Assert(newOne == nil, "Email already existed")

	userStr, err := json.Marshal(
		struct {
			FirstName, LastName, Email, Password, Phone, Address, UserName, CountryCode, CountryName string
		}{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Email:       req.Email,
			CountryCode: req.CountryCode,
			CountryName: req.CountryName,
			Password:    utility.PasswordEncrypt(req.Password),
			Phone:       req.Phone,
			Address:     req.Address,
			UserName:    req.UserName,
		},
	)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	_, err = g.Redis().Set(ctx, req.Email, userStr)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	_, err = g.Redis().Expire(ctx, req.Email, 3*60)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	verificationCode := utility.GenerateRandomCode(6)
	fmt.Printf("verification %s", verificationCode)
	_, err = g.Redis().Set(ctx, req.Email+"-verify", verificationCode)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	_, err = g.Redis().Expire(ctx, req.Email+"-verify", 3*60)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}

	err = email.SendTemplateEmail(ctx, 0, req.Email, "", email.TemplateUserRegistrationCodeVerify, "", &email.TemplateVariable{
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	if err != nil {
		return nil, err
	}

	return &auth.RegisterRes{}, nil
}
