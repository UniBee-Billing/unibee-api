package util

import (
	"context"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
)

func GetOverseaPayChannel(ctx context.Context, id int64) (channel *entity.MerchantGateway) {
	return query.GetPayChannelById(ctx, id)
}
