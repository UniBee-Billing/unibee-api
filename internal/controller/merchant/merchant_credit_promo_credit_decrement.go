package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/credit/admin"
	"unibee/utility"

	"unibee/api/merchant/credit"
)

func (c *ControllerCredit) PromoCreditDecrement(ctx context.Context, req *credit.PromoCreditDecrementReq) (res *credit.PromoCreditDecrementRes, err error) {
	utility.Assert(req.UserId > 0, "Invalid UserId")
	utility.Assert(req.Amount > 0, "Invalid Amount")
	var adminMemberId uint64 = 0
	if _interface.Context().Get(ctx).IsAdminPortalCall {
		adminMemberId = _interface.Context().Get(ctx).MerchantMember.Id
	}
	change, err := admin.CreditAccountAdminChange(ctx, &admin.CreditAccountAdminChangeInternalReq{
		UserId:        req.UserId,
		MerchantId:    _interface.GetMerchantId(ctx),
		CreditType:    consts.CreditAccountTypePromo,
		Amount:        -int64(req.Amount),
		Currency:      req.Currency,
		Name:          req.Name,
		Description:   req.Description,
		AdminMemberId: adminMemberId,
	})
	if err != nil {
		return nil, err
	}
	return &credit.PromoCreditDecrementRes{UserPromoCreditAccount: change.CreditAccount}, nil
}
