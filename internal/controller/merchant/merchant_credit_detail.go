package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/credit"
)

func (c *ControllerCredit) Detail(ctx context.Context, req *credit.DetailReq) (res *credit.DetailRes, err error) {
	utility.Assert(req.Id > 0, "Invalid Credit Account Id")
	one := query.GetCreditAccountById(ctx, req.Id)
	if one != nil {
		utility.Assert(one.MerchantId == _interface.GetMerchantId(ctx), "No permission")
	} else {
		return &credit.DetailRes{}, nil
	}
	var transactions []*bean.CreditTransaction
	_ = dao.CreditTransaction.Ctx(ctx).
		Where(dao.CreditTransaction.Columns().CreditId, req.Id).
		OmitEmpty().Scan(&transactions, true)
	return &credit.DetailRes{
		CreditAccount:      detail.ConvertToCreditAccountDetail(ctx, one),
		CreditTransactions: transactions,
	}, nil
}
