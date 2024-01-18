package balance

import (
	"context"
	"go-oversea-pay/internal/logic/gateway"
	"go-oversea-pay/internal/logic/gateway/ro"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func UserBalanceDetailQuery(ctx context.Context, merchantId int64, userId int64, channelId int64) (*ro.ChannelUserBalanceQueryInternalResp, error) {
	user := query.GetUserAccountById(ctx, uint64(userId))
	merchant := query.GetMerchantInfoById(ctx, merchantId)
	payChannel := query.GetPayChannelById(ctx, channelId)
	utility.Assert(user != nil, "user not found")
	utility.Assert(merchant != nil, "merchant not found")

	var channelUserId string
	channelUser := query.GetUserChannel(ctx, userId, channelId)
	utility.Assert(channelUser != nil, "channel User not found")
	channelUserId = channelUser.ChannelUserId

	queryResult, err := gateway.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelUserBalancesQuery(ctx, payChannel, channelUserId)
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}

func MerchantBalanceDetailQuery(ctx context.Context, merchantId int64, channelId int64) (*ro.ChannelMerchantBalanceQueryInternalResp, error) {
	merchant := query.GetMerchantInfoById(ctx, merchantId)
	payChannel := query.GetPayChannelById(ctx, channelId) // todo mark 根据 MerchantId 配置 PayChannel
	utility.Assert(merchant != nil, "merchant not found")

	queryResult, err := gateway.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelMerchantBalancesQuery(ctx, payChannel)
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}
