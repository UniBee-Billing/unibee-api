package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
	"unibee/internal/logic/jwt"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/auth"
)

func (c *ControllerAuth) LoginOtpVerify(ctx context.Context, req *auth.LoginOtpVerifyReq) (res *auth.LoginOtpVerifyRes, err error) {
	verificationCode, err := g.Redis().Get(ctx, req.Email+"-MerchantAuth-Verify")
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	utility.Assert(verificationCode != nil, "code expired")
	utility.Assert((verificationCode.String()) == req.VerificationCode, "code not match")

	var newOne *entity.MerchantMember
	newOne = query.GetMerchantMemberByEmail(ctx, req.Email)
	utility.Assert(newOne != nil, "Login Failed")

	token, err := jwt.CreateMemberPortalToken(ctx, jwt.TOKENTYPEMERCHANTMember, newOne.MerchantId, newOne.Id, req.Email)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("MerchantMember#%d", newOne.Id)), "Cache Error")
	jwt.AppendRequestCookieWithToken(ctx, token)
	return &auth.LoginOtpVerifyRes{MerchantMember: detail.ConvertMemberToDetail(ctx, newOne), Token: token}, nil
}
