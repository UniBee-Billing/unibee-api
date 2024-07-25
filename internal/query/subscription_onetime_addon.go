package query

import (
	"context"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetSubscriptionOnetimeAddonById(ctx context.Context, id uint64) (one *entity.SubscriptionOnetimeAddon) {
	if id <= 0 {
		return nil
	}
	err := dao.SubscriptionOnetimeAddon.Ctx(ctx).Where(dao.SubscriptionOnetimeAddon.Columns().Id, id).Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
