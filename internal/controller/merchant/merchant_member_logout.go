package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/jwt"
	"unibee/utility"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) Logout(ctx context.Context, req *member.LogoutReq) (res *member.LogoutRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "Merchant Member Not Found")
	utility.Assert(len(_interface.Context().Get(ctx).MerchantMember.Token) > 0, "Merchant Token Not Found")
	jwt.DelAuthToken(ctx, _interface.Context().Get(ctx).MerchantMember.Token)
	return &member.LogoutRes{}, nil
}
