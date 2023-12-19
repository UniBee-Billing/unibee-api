package query

import (
	"context"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

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
	m := dao.OverseaPayChannel.Ctx(ctx)
	err := m.Where(entity.OverseaPayChannel{ChannelType: consts.PayChannelTypeSubscription}).
		OmitEmpty().Scan(&list)
	if err != nil {
		list = nil
	}
	return
}
