package merchant

import (
	"context"
	"encoding/json"
	"fmt"
	"unibee-api/api/merchant/auth"
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
	var newOne *entity.MerchantUserAccount
	newOne = query.GetMerchantAccountByEmail(ctx, req.Email)
	utility.Assert(newOne == nil, "Email already existed")

	redisKey := fmt.Sprintf("MerchantAuth-Regist-Email:%s", req.Email)

	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}

	userStr, err := json.Marshal(
		struct {
			FirstName, LastName, Email, Password, Phone, Address, UserName string
			MerchantId                                                     uint64
		}{
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Email:      req.Email,
			Password:   utility.PasswordEncrypt(req.Password),
			Phone:      req.Phone,
			MerchantId: req.MerchantId,
			// Address:   req.Address,
			UserName: req.UserName,
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
	fmt.Println("verification ", verificationCode)
	// add merchant-verify, user-verify
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
