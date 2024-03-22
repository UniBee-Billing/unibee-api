package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/internal/logic/merchant"
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

	_, member, err := merchant.CreateMerchant(ctx, createMerchantReq)

	return &auth.RegisterVerifyRes{MerchantMember: bean.SimplifyMerchantMember(member)}, nil
}
