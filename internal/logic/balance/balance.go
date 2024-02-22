package balance

import (
	"context"
	"unibee-api/internal/logic/gateway/api"
	"unibee-api/internal/logic/gateway/ro"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func UserBalanceDetailQuery(ctx context.Context, merchantId uint64, userId int64, gatewayId int64) (*ro.GatewayUserDetailQueryInternalResp, error) {
	user := query.GetUserAccountById(ctx, uint64(userId))
	merchant := query.GetMerchantInfoById(ctx, merchantId)
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(user != nil, "user not found")
	utility.Assert(merchant != nil, "merchant not found")

	queryResult, err := api.GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayUserDetailQuery(ctx, gateway, userId)
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}

func MerchantBalanceDetailQuery(ctx context.Context, merchantId uint64, gatewayId int64) (*ro.GatewayMerchantBalanceQueryInternalResp, error) {
	merchant := query.GetMerchantInfoById(ctx, merchantId)
	gateway := query.GetGatewayById(ctx, gatewayId) // todo mark 根据 MerchantId 配置 Gateway
	utility.Assert(merchant != nil, "merchant not found")

	queryResult, err := api.GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayMerchantBalancesQuery(ctx, gateway)
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}
