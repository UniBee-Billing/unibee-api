package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	balance2 "unibee/internal/logic/balance"

	"unibee/api/merchant/balance"
)

func (c *ControllerBalance) DetailQuery(ctx context.Context, req *balance.DetailQueryReq) (res *balance.DetailQueryRes, err error) {
	balanceResult, err := balance2.MerchantBalanceDetailQuery(ctx, _interface.GetMerchantId(ctx), req.GatewayId)

	return &balance.DetailQueryRes{
		AvailableBalance:       balanceResult.AvailableBalance,
		ConnectReservedBalance: balanceResult.ConnectReservedBalance,
		PendingBalance:         balanceResult.PendingBalance,
	}, nil
}
