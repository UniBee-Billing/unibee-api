package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	member2 "unibee/internal/logic/member"
	"unibee/utility"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) PasswordReset(ctx context.Context, req *member.PasswordResetReq) (res *member.PasswordResetRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "Merchant User Not Found")
	utility.Assert(len(_interface.Context().Get(ctx).MerchantMember.Token) > 0, "Merchant Token Not Found")

	member2.ChangeMerchantMemberPassword(ctx, _interface.Context().Get(ctx).MerchantMember.Email, req.OldPassword, req.NewPassword)
	return &member.PasswordResetRes{}, nil
}
