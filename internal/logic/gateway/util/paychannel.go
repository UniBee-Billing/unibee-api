package util

import (
	"context"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
)

func GetGatewayById(ctx context.Context, id int64) (gateway *entity.MerchantGateway) {
	return query.GetGatewayById(ctx, id)
}
