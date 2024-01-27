package out

import (
	"context"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func queryAndCreateChannelUserId(ctx context.Context, payChannel *entity.OverseaPayChannel, userId int64) string {
	// todo mark sync lock control
	channelUser := query.GetUserChannel(ctx, userId, int64(payChannel.Id))
	if channelUser == nil {
		user := query.GetUserAccountById(ctx, uint64(userId))
		utility.Assert(user != nil, "user not found")
		utility.Assert(len(user.Email) > 0, "invalid user email")
		create, err := GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelUserCreate(ctx, payChannel, user)
		utility.Assert(err == nil, "queryAndCreateChannelUser:"+err.Error())
		_, err = query.SaveUserChannel(ctx, userId, int64(payChannel.Id), create.ChannelUserId)
		utility.Assert(err == nil, "queryAndCreateChannelUser:"+err.Error())
		return create.ChannelUserId
	} else {
		return channelUser.ChannelUserId
	}
}
