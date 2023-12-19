package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetSubscriptionPlanById(ctx context.Context, id int64) (one *entity.SubscriptionPlan) {
	err := dao.SubscriptionPlan.Ctx(ctx).Where(entity.SubscriptionPlan{Id: uint64(id)}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
