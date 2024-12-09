package merchant

import (
	"context"
	"strings"
	_interface "unibee/internal/interface"
	credit2 "unibee/internal/logic/credit/transaction"

	"unibee/api/merchant/credit"
)

func (c *ControllerCredit) CreditTransactionList(ctx context.Context, req *credit.CreditTransactionListReq) (res *credit.CreditTransactionListRes, err error) {
	list, err := credit2.CreditTransactionList(ctx, &credit2.CreditTransactionListInternalReq{
		MerchantId:       _interface.GetMerchantId(ctx),
		AccountType:      req.AccountType,
		UserId:           req.UserId,
		Email:            strings.ToLower(req.Email),
		Currency:         req.Currency,
		SortField:        req.SortField,
		SortType:         req.SortType,
		TransactionTypes: req.TransactionTypes,
		Page:             req.Page,
		Count:            req.Count,
		CreateTimeStart:  req.CreateTimeStart,
		CreateTimeEnd:    req.CreateTimeEnd,
	})
	if err != nil {
		return nil, err
	}
	return &credit.CreditTransactionListRes{
		CreditTransactions: list.CreditTransactions,
		Total:              list.Total,
	}, nil
}
