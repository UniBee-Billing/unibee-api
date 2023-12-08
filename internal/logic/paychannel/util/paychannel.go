package util

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetOverseaPayChannel(ctx context.Context, id uint64) (channel *entity.OverseaPayChannel) {
	var (
		one *entity.OverseaPayChannel
	)
	err := dao.OverseaPayChannel.Ctx(ctx).Where(entity.OverseaPayChannel{Id: id}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}
