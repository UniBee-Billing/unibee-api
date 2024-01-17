package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetSubscriptionVatRateChannel(ctx context.Context, vatRateId uint64, channelId uint64) (one *entity.SubscriptionVatRateChannel) {
	err := dao.SubscriptionVatRateChannel.Ctx(ctx).Where(entity.SubscriptionVatRateChannel{VatRateId: int64(vatRateId), ChannelId: int64(channelId)}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetSubscriptionVatRateChannelById(ctx context.Context, id int64) (one *entity.SubscriptionVatRateChannel) {
	err := dao.SubscriptionVatRateChannel.Ctx(ctx).Where(entity.SubscriptionVatRateChannel{Id: uint64(id)}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
