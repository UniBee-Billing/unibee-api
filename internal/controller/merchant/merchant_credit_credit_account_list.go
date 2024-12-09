package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	credit2 "unibee/internal/logic/credit"

	"unibee/api/merchant/credit"
)

func (c *ControllerCredit) CreditAccountList(ctx context.Context, req *credit.CreditAccountListReq) (res *credit.CreditAccountListRes, err error) {
	list, err := credit2.CreditAccountList(ctx, &credit2.CreditAccountListInternalReq{
		MerchantId:      _interface.GetMerchantId(ctx),
		UserId:          req.UserId,
		Email:           req.Email,
		SortField:       req.SortField,
		SortType:        req.SortType,
		Page:            req.Page,
		Count:           req.Count,
		CreateTimeStart: req.CreateTimeStart,
		CreateTimeEnd:   req.CreateTimeEnd,
	})
	if err != nil {
		return nil, err
	}
	return &credit.CreditAccountListRes{
		CreditAccounts: list.CreditAccounts,
		Total:          list.Total,
	}, nil
}
