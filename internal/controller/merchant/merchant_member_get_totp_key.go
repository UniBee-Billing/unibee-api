package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/totp"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) GetTotpKey(ctx context.Context, req *member.GetTotpKeyReq) (res *member.GetTotpKeyRes, err error) {
	one := query.GetMerchantMemberById(ctx, _interface.Context().Get(ctx).MerchantMember.Id)
	utility.Assert(one != nil, "Merchant Member Not Found")
	utility.Assert(one.TotpValidatorType == 0, "Already setup 2FA")
	key, url, err := totp.GetMerchantMemberTotpSecret(ctx, req.TotpType, one.Email)
	utility.AssertError(err, "Get 2FA error")
	return &member.GetTotpKeyRes{
		TotpKey:        key,
		TotpResumeCode: utility.MD5(key),
		TotpUrl:        url,
		TotpType:       req.TotpType,
	}, nil
}
