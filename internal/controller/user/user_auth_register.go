package user

import (
	"context"
	"encoding/json"
	"fmt"
	"unibee/api/user/auth"
	_interface "unibee/internal/interface"
	auth2 "unibee/internal/logic/auth"
	"unibee/internal/logic/email"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"

	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

const CacheKeyUserRegisterPrefix = "CacheKeyUserRegisterPrefix-"

func (c *ControllerAuth) Register(ctx context.Context, req *auth.RegisterReq) (res *auth.RegisterRes, err error) {
	utility.Assert(len(req.Email) > 0, "Email Needed")
	utility.Assert(utility.IsEmailValid(req.Email), "Invalid Email")

	redisKey := fmt.Sprintf("UserAuth-Regist-Email:%s", req.Email)

	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}

	var newOne *entity.UserAccount
	newOne = query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email) //Id(ctx, user.Id)
	utility.Assert(newOne == nil, "Email already existed")

	userStr, err := json.Marshal(
		&auth2.NewReq{
			Email:       req.Email,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Password:    req.Password,
			Phone:       req.Phone,
			Address:     req.Address,
			UserName:    req.UserName,
			CountryCode: req.CountryCode,
			CountryName: req.CountryName,
		},
	)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Set(ctx, CacheKeyUserRegisterPrefix+req.Email, userStr)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, CacheKeyUserRegisterPrefix+req.Email, 3*60)
	utility.AssertError(err, "Server Error")
	verificationCode := utility.GenerateRandomCode(6)
	fmt.Printf("verification %s", verificationCode)
	_, err = g.Redis().Set(ctx, CacheKeyUserRegisterPrefix+req.Email+"-verify", verificationCode)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, CacheKeyUserRegisterPrefix+req.Email+"-verify", 3*60)
	utility.AssertError(err, "Server Error")

	err = email.SendTemplateEmail(ctx, _interface.GetMerchantId(ctx), req.Email, "", email.TemplateUserRegistrationCodeVerify, "", &email.TemplateVariable{
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	utility.AssertError(err, "Server Error")

	return &auth.RegisterRes{}, nil
}
