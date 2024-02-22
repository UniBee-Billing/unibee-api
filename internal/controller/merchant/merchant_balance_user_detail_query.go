package merchant

import (
	"context"
	_interface "unibee-api/internal/interface"
	balance2 "unibee-api/internal/logic/balance"

	"unibee-api/api/merchant/balance"
)

func (c *ControllerBalance) UserDetailQuery(ctx context.Context, req *balance.UserDetailQueryReq) (res *balance.UserDetailQueryRes, err error) {
	balanceResult, err := balance2.UserBalanceDetailQuery(ctx, _interface.GetMerchantId(ctx), req.UserId, req.GatewayId)

	return &balance.UserDetailQueryRes{
		Balance:              balanceResult.Balance,
		CashBalance:          balanceResult.CashBalance,
		InvoiceCreditBalance: balanceResult.InvoiceCreditBalance,
	}, nil
}
