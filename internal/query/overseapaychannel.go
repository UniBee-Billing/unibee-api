package query

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetPayChannelById(ctx context.Context, id int64) (one *entity.OverseaPayChannel) {
	m := dao.OverseaPayChannel.Ctx(ctx)
	err := m.Where(entity.OverseaPayChannel{Id: uint64(id)}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetPayChannelByChannel(ctx context.Context, channel string) (one *entity.OverseaPayChannel) {
	err := dao.OverseaPayChannel.Ctx(ctx).Where(entity.OverseaPayChannel{Channel: channel}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetPaymentTypePayChannelById(ctx context.Context, id int64) (one *entity.OverseaPayChannel) {
	m := dao.OverseaPayChannel.Ctx(ctx)
	err := m.Where(entity.OverseaPayChannel{Id: uint64(id)}).
		Where(m.Builder().
			Where(entity.OverseaPayChannel{ChannelType: consts.PayChannelTypePayment}).WhereOr("channel_type is null")).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetSubscriptionTypePayChannelById(ctx context.Context, id int64) (one *entity.OverseaPayChannel) {
	m := dao.OverseaPayChannel.Ctx(ctx)
	err := m.Where(entity.OverseaPayChannel{Id: uint64(id), ChannelType: consts.PayChannelTypeSubscription}).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetListSubscriptionTypePayChannels(ctx context.Context) (list *[]entity.OverseaPayChannel) {
	var data []entity.OverseaPayChannel
	err := dao.OverseaPayChannel.Ctx(ctx).Where(entity.OverseaPayChannel{ChannelType: consts.PayChannelTypeSubscription}).
		OmitEmpty().Scan(&data)
	if err != nil {
		g.Log().Errorf(ctx, "GetListSubscriptionTypePayChannels error:%s", err)
		return nil
	}
	return &data
}
