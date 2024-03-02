package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/auth"
	"unibee/utility"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) PasswordReset(ctx context.Context, req *member.PasswordResetReq) (res *member.PasswordResetRes, err error) {
	utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "Merchant User Not Found")
	utility.Assert(len(_interface.BizCtx().Get(ctx).MerchantMember.Token) > 0, "Merchant Token Not Found")

	auth.ChangeMerchantMemberPassword(ctx, _interface.GetMerchantId(ctx), _interface.BizCtx().Get(ctx).MerchantMember.Email, req.OldPassword, req.NewPassword)
	return &member.PasswordResetRes{}, nil
}
