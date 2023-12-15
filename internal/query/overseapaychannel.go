package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetOverseaPayChannelById(ctx context.Context, id int64) (one *entity.OverseaPayChannel) {
	err := dao.OverseaPayChannel.Ctx(ctx).Where(entity.OverseaPayChannel{Id: uint64(id)}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
