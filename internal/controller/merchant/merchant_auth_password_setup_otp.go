package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/jwt"
	auth2 "unibee/internal/logic/member"
	"unibee/internal/logic/totp/client_activity"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/auth"
)

func (c *ControllerAuth) PasswordSetupOtp(ctx context.Context, req *auth.PasswordSetupOtpReq) (res *auth.PasswordSetupOtpRes, err error) {
	merchantMember := query.GetMerchantMemberByEmail(ctx, req.Email)
	utility.Assert(merchantMember != nil, "User Not Found")
	utility.Assert(merchantMember.TotpValidatorSecret == req.SetupToken, "Invalid setupToken")
	auth2.ChangeMerchantMemberPasswordWithOutOldVerify(ctx, req.Email, req.NewPassword)
	_, err = dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().TotpValidatorSecret: "",
		dao.MerchantMember.Columns().GmtModify:           gtime.Now(),
	}).Where(dao.MerchantMember.Columns().Id, merchantMember.Id).Update()
	utility.AssertError(err, "setup failed")
	// 9. Generate system JWT token
	token, err := jwt.CreateMemberPortalToken(ctx, jwt.TOKENTYPEMERCHANTMember, merchantMember.MerchantId, merchantMember.Id, merchantMember.Email)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("MerchantMember#%d", merchantMember.Id)), "Cache Error")

	// 10. Set cookies
	g.RequestFromCtx(ctx).Cookie.Set(jwt.MERCHANT_TYPE_TOKEN_COOKIE_KEY, token)
	jwt.AppendRequestCookieWithToken(ctx, token)
	client_activity.UpdateClientIdentityLoginTime(ctx, _interface.Context().Get(ctx).ClientIdentity)
	return &auth.PasswordSetupOtpRes{MerchantMember: detail.ConvertMemberToDetail(ctx, merchantMember), Token: token}, nil
}
