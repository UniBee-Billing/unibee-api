package merchant

import (
	"context"
	"encoding/json"
	"fmt"
	"unibee/api/merchant/auth"
	"unibee/internal/cmd/config"
	"unibee/internal/cmd/i18n"
	"unibee/internal/logic/merchant"
	"unibee/internal/logic/middleware"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"

	entity "unibee/internal/model/entity/default"
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
		utility.Assert(false, i18n.LocalizationFormat(ctx, "{#ClickTooFast}"))
	}
	internalReq := &merchant.CreateMerchantInternalReq{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		Phone:     req.Phone,
		UserName:  req.UserName,
	}
	userStr, err := json.Marshal(
		internalReq,
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

	list := query.GetActiveMerchantList(ctx)
	if len(list) >= 2 {
		utility.Assert(config.GetConfigInstance().Mode == "cloud", "Register multi merchants should contain valid mode")
		var containPremiumMerchant = false
		for _, one := range list {
			if middleware.IsPremiumVersion(ctx, one.Id) {
				containPremiumMerchant = true
				break
			}
		}
		utility.Assert(containPremiumMerchant, "Feature register multi merchants need premium license, contact us directly if needed")
	}

	merchant.SendMerchantRegisterEmail(ctx, internalReq, verificationCode)
	return &auth.RegisterRes{}, nil
}
