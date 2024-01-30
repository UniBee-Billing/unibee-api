package balance

import (
	"context"
	"go-oversea-pay/internal/logic/channel/out"
	"go-oversea-pay/internal/logic/channel/ro"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func UserBalanceDetailQuery(ctx context.Context, merchantId int64, userId int64, channelId int64) (*ro.ChannelUserDetailQueryInternalResp, error) {
	user := query.GetUserAccountById(ctx, uint64(userId))
	merchant := query.GetMerchantInfoById(ctx, merchantId)
	payChannel := query.GetPayChannelById(ctx, channelId)
	utility.Assert(user != nil, "user not found")
	utility.Assert(merchant != nil, "merchant not found")

	queryResult, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelUserDetailQuery(ctx, payChannel, userId)
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}

func MerchantBalanceDetailQuery(ctx context.Context, merchantId int64, channelId int64) (*ro.ChannelMerchantBalanceQueryInternalResp, error) {
	merchant := query.GetMerchantInfoById(ctx, merchantId)
	payChannel := query.GetPayChannelById(ctx, channelId) // todo mark 根据 MerchantId 配置 PayChannel
	utility.Assert(merchant != nil, "merchant not found")

	queryResult, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelMerchantBalancesQuery(ctx, payChannel)
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}
