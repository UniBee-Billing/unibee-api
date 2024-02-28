package query

import (
	"context"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func GetSubscriptionVatRateChannel(ctx context.Context, vatRateId uint64, gatewayId uint64) (one *entity.GatewayVatRate) {
	if gatewayId <= 0 || vatRateId <= 0 {
		return nil
	}
	err := dao.GatewayVatRate.Ctx(ctx).Where(entity.GatewayVatRate{VatRateId: int64(vatRateId), GatewayId: gatewayId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetSubscriptionVatRateChannelById(ctx context.Context, id int64) (one *entity.GatewayVatRate) {
	if id <= 0 {
		return nil
	}
	err := dao.GatewayVatRate.Ctx(ctx).Where(entity.GatewayVatRate{Id: uint64(id)}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
