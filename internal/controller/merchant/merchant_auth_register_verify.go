package merchant

import (
	"context"
	"unibee/api/bean/detail"
	"unibee/internal/cmd/config"
	"unibee/internal/logic/merchant"
	"unibee/internal/logic/middleware"
	"unibee/internal/logic/vat_gateway/setup"
	"unibee/internal/query"
	"unibee/utility"

	"encoding/json"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/merchant/auth"
)

func (c *ControllerAuth) RegisterVerify(ctx context.Context, req *auth.RegisterVerifyReq) (res *auth.RegisterVerifyRes, err error) {
	verificationCode, err := g.Redis().Get(ctx, CacheKeyMerchantRegisterPrefix+req.Email+"-verify")
	utility.AssertError(err, "Server Error")
	utility.Assert(verificationCode != nil, "Invalid Code")
	utility.Assert((verificationCode.String()) == req.VerificationCode, "Invalid Code")
	userStr, err := g.Redis().Get(ctx, CacheKeyMerchantRegisterPrefix+req.Email)
	utility.AssertError(err, "Server Error")
	utility.Assert(userStr != nil, "Invalid Code")
	var createMerchantReq *merchant.CreateMerchantInternalReq
	err = json.Unmarshal([]byte(userStr.String()), &createMerchantReq)

	list := query.GetActiveMerchantList(ctx)
	if len(list) > 2 {
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

	_, member, err := merchant.CreateMerchant(ctx, createMerchantReq)
	if member != nil {
		_ = setup.InitMerchantDefaultVatGateway(ctx, member.MerchantId)
	}
	return &auth.RegisterVerifyRes{MerchantMember: detail.ConvertMemberToDetail(ctx, member)}, nil
}
