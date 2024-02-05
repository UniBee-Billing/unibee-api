package merchant

import (
	"context"
	balance2 "go-oversea-pay/internal/logic/balance"

	"go-oversea-pay/api/merchant/balance"
)

func (c *ControllerBalance) DetailQuery(ctx context.Context, req *balance.DetailQueryReq) (res *balance.DetailQueryRes, err error) {
	balanceResult, err := balance2.MerchantBalanceDetailQuery(ctx, req.MerchantId, req.GatewayId)

	return &balance.DetailQueryRes{
		AvailableBalance:       balanceResult.AvailableBalance,
		ConnectReservedBalance: balanceResult.ConnectReservedBalance,
		PendingBalance:         balanceResult.PendingBalance,
	}, nil
}
