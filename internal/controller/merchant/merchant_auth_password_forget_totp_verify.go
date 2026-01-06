package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	auth2 "unibee/internal/logic/member"
	"unibee/internal/logic/totp"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/auth"
)

func (c *ControllerAuth) PasswordForgetTotpVerify(ctx context.Context, req *auth.PasswordForgetTotpVerifyReq) (res *auth.PasswordForgetTotpVerifyRes, err error) {
	var one *entity.MerchantMember
	one = query.GetMerchantMemberByEmail(ctx, req.Email)
	utility.Assert(one != nil, "Member Not Found")

	//verificationCode, err := g.Redis().Get(ctx, req.Email+"-MerchantAuth-PasswordForgetOtp-Verify")
	//if err != nil {
	//	return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	//}
	//utility.Assert(verificationCode != nil, "code expired")
	//utility.Assert((verificationCode.String()) == req.VerificationCode, "code not match")

	if one.TotpValidatorType > 0 && len(one.TotpValidatorSecret) > 0 {
		utility.Assert(len(req.TotpCode) > 0, "2FA Code Needed")
		utility.Assert(totp.ValidateMerchantMemberTotp(ctx, one.TotpValidatorType, one.Email, one.TotpValidatorSecret, req.TotpCode, _interface.Context().Get(ctx).ClientIdentity), "Invalid 2FA Code")
	} else {
		utility.Assert(false, "2FA need setup first")
	}

	auth2.ChangeMerchantMemberPasswordWithOutOldVerify(ctx, req.Email, req.NewPassword)
	return &auth.PasswordForgetTotpVerifyRes{}, nil
}
