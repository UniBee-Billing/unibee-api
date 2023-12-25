package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetSubscriptionById(ctx context.Context, id int64) (one *entity.Subscription) {
	err := dao.Subscription.Ctx(ctx).Where(entity.Subscription{Id: uint64(id)}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
