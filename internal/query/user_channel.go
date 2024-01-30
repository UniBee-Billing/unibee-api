package query

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
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

func GetUserChannelByChannelUserId(ctx context.Context, channelUserId string, channelId int64) (one *entity.SubscriptionUserChannel) {
	utility.Assert(len(channelUserId) > 0, "invalid channelUserId")
	utility.Assert(channelId > 0, "invalid channelId")
	err := dao.SubscriptionUserChannel.Ctx(ctx).Where(entity.SubscriptionUserChannel{ChannelUserId: channelUserId, ChannelId: channelId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func CreateOrUpdateChannelUser(ctx context.Context, userId int64, channelId int64, channelUserId string, channelDefaultPaymentMethod string) (*entity.SubscriptionUserChannel, error) {
	utility.Assert(userId > 0, "invalid userId")
	utility.Assert(channelId > 0, "invalid channelId")
	utility.Assert(len(channelUserId) > 0, "invalid channelUserId")
	one := GetUserChannel(ctx, userId, channelId)
	if one == nil {
		one = &entity.SubscriptionUserChannel{
			UserId:                      userId,
			ChannelId:                   channelId,
			ChannelUserId:               channelUserId,
			ChannelDefaultPaymentMethod: channelDefaultPaymentMethod,
		}
		result, err := dao.SubscriptionUserChannel.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateChannelUser record insert failure %s`, err)
			return nil, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		one.Id = uint64(uint(id))
	} else {
		one.ChannelDefaultPaymentMethod = channelDefaultPaymentMethod
		_, err := dao.SubscriptionUserChannel.Ctx(ctx).Data(g.Map{
			dao.SubscriptionUserChannel.Columns().ChannelDefaultPaymentMethod: channelDefaultPaymentMethod,
		}).Where(dao.SubscriptionUserChannel.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateChannelUser record insert failure %s`, err)
			return nil, err
		}
	}
	return one, nil
}
