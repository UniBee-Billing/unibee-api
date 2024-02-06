package util

import (
	"context"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
)

func GetGatewayById(ctx context.Context, id int64) (gateway *entity.MerchantGateway) {
	return query.GetGatewayById(ctx, id)
}
