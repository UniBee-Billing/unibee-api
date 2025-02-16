package balance

import (
	"context"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/query"
	"unibee/utility"
)

func UserBalanceDetailQuery(ctx context.Context, merchantId uint64, userId uint64, gatewayId uint64) (*gateway_bean.GatewayUserDetailQueryResp, error) {
	user := query.GetUserAccountById(ctx, uint64(userId))
	merchant := query.GetMerchantById(ctx, merchantId)
	gateway := query.GetGatewayById(ctx, uint64(gatewayId))
	utility.Assert(user != nil, "user not found")
	utility.Assert(merchant != nil, "merchant not found")

	queryResult, err := api.GetGatewayServiceProvider(ctx, gateway.Id).GatewayUserDetailQuery(ctx, gateway, userId)
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}

func MerchantBalanceDetailQuery(ctx context.Context, merchantId uint64, gatewayId uint64) (*gateway_bean.GatewayMerchantBalanceQueryResp, error) {
	merchant := query.GetMerchantById(ctx, merchantId)
	gateway := query.GetGatewayById(ctx, uint64(gatewayId))
	utility.Assert(merchant != nil, "merchant not found")

	queryResult, err := api.GetGatewayServiceProvider(ctx, gateway.Id).GatewayMerchantBalancesQuery(ctx, gateway)
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}
