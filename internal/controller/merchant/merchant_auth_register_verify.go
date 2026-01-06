package merchant

import (
	"context"
	"fmt"
	"unibee/api/bean/detail"
	"unibee/internal/cmd/config"
	"unibee/internal/logic/jwt"
	"unibee/internal/logic/merchant"
	"unibee/internal/logic/middleware/license"
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
			if license.IsPremiumVersion(ctx, one.Id) {
				containPremiumMerchant = true
				break
			}
		}
		utility.Assert(containPremiumMerchant, "Feature register multi merchants need premium license, contact us directly if needed")
	}

	_, member, err := merchant.CreateMerchant(ctx, createMerchantReq)
	utility.AssertError(err, "CreateMerchant Error")

	// 9. Generate system JWT token
	token, err := jwt.CreateMemberPortalToken(ctx, jwt.TOKENTYPEMERCHANTMember, member.MerchantId, member.Id, member.Email)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("MerchantMember#%d", member.Id)), "Cache Error")

	// 10. Set cookies
	g.RequestFromCtx(ctx).Cookie.Set(jwt.MERCHANT_TYPE_TOKEN_COOKIE_KEY, token)
	jwt.AppendRequestCookieWithToken(ctx, token)
	return &auth.RegisterVerifyRes{MerchantMember: detail.ConvertMemberToDetail(ctx, member), Token: token}, nil
}
