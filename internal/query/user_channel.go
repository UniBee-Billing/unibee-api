package query

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
)

func GetUserChannel(ctx context.Context, userId int64, channelId int64) (one *entity.SubscriptionUserChannel) {
	utility.Assert(userId > 0, "invalid userId")
	utility.Assert(channelId > 0, "invalid channelId")
	err := dao.SubscriptionUserChannel.Ctx(ctx).Where(entity.SubscriptionUserChannel{UserId: userId, ChannelId: channelId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func SaveUserChannel(ctx context.Context, userId int64, channelId int64, channelUserId string) (*entity.SubscriptionUserChannel, error) {
	utility.Assert(userId > 0, "invalid userId")
	utility.Assert(channelId > 0, "invalid channelId")
	utility.Assert(len(channelUserId) > 0, "invalid channelUserId")
	one := &entity.SubscriptionUserChannel{
		UserId:        userId,
		ChannelId:     channelId,
		ChannelUserId: channelUserId,
	}
	result, err := dao.SubscriptionUserChannel.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionCreate record insert failure %s`, err)
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	one.Id = uint64(uint(id))
	return one, nil
}
