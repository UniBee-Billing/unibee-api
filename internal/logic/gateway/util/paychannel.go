package util

import (
	"context"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

func GetGatewayById(ctx context.Context, id int64) (gateway *entity.MerchantGateway) {
	return query.GetGatewayById(ctx, id)
}
