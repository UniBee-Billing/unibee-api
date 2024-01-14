package util

import (
	"context"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
)

func GetOverseaPayChannel(ctx context.Context, id int64) (channel *entity.OverseaPayChannel) {
	return query.GetPayChannelById(ctx, id)
}
