package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetSubscriptionVatRateChannel(ctx context.Context, vatRateId uint64, channelId uint64) (one *entity.ChannelVatRate) {
	if channelId <= 0 || vatRateId <= 0 {
		return nil
	}
	err := dao.ChannelVatRate.Ctx(ctx).Where(entity.ChannelVatRate{VatRateId: int64(vatRateId), ChannelId: int64(channelId)}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetSubscriptionVatRateChannelById(ctx context.Context, id int64) (one *entity.ChannelVatRate) {
	if id <= 0 {
		return nil
	}
	err := dao.ChannelVatRate.Ctx(ctx).Where(entity.ChannelVatRate{Id: uint64(id)}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
