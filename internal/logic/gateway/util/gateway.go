package util

import (
	"context"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

func GetGatewayById(ctx context.Context, id uint64) (gateway *entity.MerchantGateway) {
	return query.GetGatewayById(ctx, id)
}
