package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/credit/admin"
	"unibee/utility"

	"unibee/api/merchant/credit"
)

func (c *ControllerCredit) PromoCreditIncrement(ctx context.Context, req *credit.PromoCreditIncrementReq) (res *credit.PromoCreditIncrementRes, err error) {
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
		Amount:        int64(req.Amount),
		Currency:      req.Currency,
		Name:          req.Name,
		Description:   req.Description,
		AdminMemberId: adminMemberId,
	})
	if err != nil {
		return nil, err
	}
	return &credit.PromoCreditIncrementRes{UserPromoCreditAccount: change.CreditAccount}, nil
}
