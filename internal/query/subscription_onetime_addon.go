package query

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/plan/period"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
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

func GetSubscriptionOnetimeAddons(ctx context.Context, sub *entity.Subscription) (list []*entity.SubscriptionOnetimeAddon) {
	if sub == nil {
		return make([]*entity.SubscriptionOnetimeAddon, 0)
	}
	if sub.Status != consts.SubStatusActive && sub.Status != consts.SubStatusIncomplete {
		return make([]*entity.SubscriptionOnetimeAddon, 0)
	}
	timeNow := utility.MaxInt64(gtime.Timestamp(), sub.TestClock)
	query := dao.SubscriptionOnetimeAddon.Ctx(ctx)
	_ = query.
		Where(dao.SubscriptionOnetimeAddon.Columns().UserId, sub.UserId).
		Where(dao.SubscriptionOnetimeAddon.Columns().SubscriptionId, sub.SubscriptionId).
		Where(dao.SubscriptionOnetimeAddon.Columns().Status, 2).
		Where(query.Builder().WhereOrGTE(dao.SubscriptionOnetimeAddon.Columns().PeriodEnd, timeNow).
			WhereOr(dao.SubscriptionOnetimeAddon.Columns().PeriodEnd, 0)).
		Scan(&list)
	result := make([]*entity.SubscriptionOnetimeAddon, 0)
	for _, one := range list {
		if one.PeriodEnd == 0 && one.GmtModify != nil && period.GetPeriodEndFromStart(ctx, one.GmtModify.Timestamp(), one.GmtModify.Timestamp(), one.AddonId) > timeNow {
			//addon := GetPlanById(ctx, one.AddonId)
			//if addon != nil {
			//	periodEnd :=
			//}
			result = append(result, one)
		} else if one.PeriodEnd != 0 {
			result = append(result, one)
		}
	}
	return result
}
