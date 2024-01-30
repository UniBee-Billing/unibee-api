package out

import (
	"context"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func queryAndCreateChannelUserId(ctx context.Context, payChannel *entity.MerchantChannelConfig, userId int64) *entity.ChannelUser {
	// todo mark sync lock control
	channelUser := query.GetUserChannel(ctx, userId, int64(payChannel.Id))
	if channelUser == nil {
		user := query.GetUserAccountById(ctx, uint64(userId))
		utility.Assert(user != nil, "user not found")
		utility.Assert(len(user.Email) > 0, "invalid user email")
		create, err := GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelUserCreate(ctx, payChannel, user)
		utility.AssertError(err, "DoRemoteChannelUserCreate")
		channelUser, err = query.CreateOrUpdateChannelUser(ctx, userId, int64(payChannel.Id), create.ChannelUserId, "")
		utility.AssertError(err, "CreateOrUpdateChannelUser")
		return channelUser
	} else {
		if len(channelUser.ChannelDefaultPaymentMethod) == 0 {
			//no default payment method, query it
			detailQuery, err := GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelUserDetailQuery(ctx, payChannel, channelUser.UserId)
			utility.AssertError(err, "DoRemoteChannelUserDetailQuery")
			if len(detailQuery.DefaultPaymentMethod) > 0 {
				channelUser, err = query.CreateOrUpdateChannelUser(ctx, userId, int64(payChannel.Id), channelUser.ChannelUserId, detailQuery.DefaultPaymentMethod)
				channelUser.ChannelDefaultPaymentMethod = detailQuery.DefaultPaymentMethod
			}
		}
		return channelUser
	}
}
