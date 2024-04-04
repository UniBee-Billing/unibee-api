package merchant

import (
	"context"
	"encoding/json"
	"fmt"
	"unibee/api/merchant/auth"
	"unibee/internal/cmd/config"
	"unibee/internal/logic/merchant"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"

	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

const CacheKeyMerchantRegisterPrefix = "CacheKeyMerchantRegisterPrefix-"

func (c *ControllerAuth) Register(ctx context.Context, req *auth.RegisterReq) (res *auth.RegisterRes, err error) {
	utility.Assert(len(req.Email) > 0, "Email Needed")
	utility.Assert(utility.IsEmailValid(req.Email), "Invalid Email")
	var newOne *entity.MerchantMember
	newOne = query.GetMerchantMemberByEmail(ctx, req.Email)
	utility.Assert(newOne == nil, "Email already existed")
	redisKey := fmt.Sprintf("MerchantAuth-Regist-Email:%s", req.Email)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}
	utility.Assert(config.GetConfigInstance().Mode == "cloud", "unsupported")

	userStr, err := json.Marshal(
		&merchant.CreateMerchantInternalReq{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  req.Password,
			Phone:     req.Phone,
			UserName:  req.UserName,
		},
	)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Set(ctx, CacheKeyMerchantRegisterPrefix+req.Email, userStr)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, CacheKeyMerchantRegisterPrefix+req.Email, 3*60)
	utility.AssertError(err, "Server Error")
	verificationCode := utility.GenerateRandomCode(6)
	fmt.Println("verification ", verificationCode)
	_, err = g.Redis().Set(ctx, CacheKeyMerchantRegisterPrefix+req.Email+"-verify", verificationCode)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, CacheKeyMerchantRegisterPrefix+req.Email+"-verify", 3*60)
	utility.AssertError(err, "Server Error")

	merchant.SendMerchantRegisterEmail(ctx, req.Email, verificationCode)
	return &auth.RegisterRes{}, nil
}
