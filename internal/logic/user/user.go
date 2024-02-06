package user

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/query"
)

func UpdateUserDefaultSubscription(ctx context.Context, userId int64, subscriptionId string) {
	if userId > 0 && len(subscriptionId) > 0 {
		one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
		if one != nil {
			_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().SubscriptionId: subscriptionId,
				dao.UserAccount.Columns().GmtModify:      gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, userId).OmitNil().Update()
			if err != nil {
				g.Log().Errorf(ctx, "UpdateUserDefaultSubscription err:%s", err.Error())
			}
		}
	}
}

func UpdateUserDefaultVatNumber(ctx context.Context, userId int64, vatNumber string) {
	if userId > 0 && len(vatNumber) > 0 {
		_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
			dao.UserAccount.Columns().VATNumber: vatNumber,
			dao.UserAccount.Columns().GmtModify: gtime.Now(),
		}).Where(dao.UserAccount.Columns().Id, userId).OmitNil().Update()
		if err != nil {
			g.Log().Errorf(ctx, "UpdateUserDefaultVatNumber err:%s", err.Error())
		}
	}
}
