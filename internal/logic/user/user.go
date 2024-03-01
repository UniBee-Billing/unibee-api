package user

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/query"
)

func UpdateUserDefaultSubscriptionForUpdate(ctx context.Context, userId int64, subscriptionId string) {
	if userId > 0 && len(subscriptionId) > 0 {
		one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
		user := query.GetUserAccountById(ctx, uint64(userId))
		if one != nil && user != nil && user.SubscriptionId == subscriptionId {
			_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().SubscriptionId:     subscriptionId,
				dao.UserAccount.Columns().SubscriptionStatus: one.Status,
				dao.UserAccount.Columns().GmtModify:          gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, userId).OmitNil().Update()
			if err != nil {
				g.Log().Errorf(ctx, "UpdateUserDefaultSubscriptionForPaymentSuccess err:%s", err.Error())
			}
		}
	}
}

func UpdateUserDefaultSubscriptionForPaymentSuccess(ctx context.Context, userId int64, subscriptionId string) {
	if userId > 0 && len(subscriptionId) > 0 {
		one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
		user := query.GetUserAccountById(ctx, uint64(userId))
		var subName = ""
		if one != nil && user != nil {
			plan := query.GetPlanById(ctx, one.PlanId)
			if plan != nil {
				subName = plan.PlanName
			}
			_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().SubscriptionId:     subscriptionId,
				dao.UserAccount.Columns().SubscriptionStatus: one.Status,
				dao.UserAccount.Columns().SubscriptionName:   subName,
				dao.UserAccount.Columns().RecurringAmount:    user.RecurringAmount + one.Amount,
				dao.UserAccount.Columns().GmtModify:          gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, userId).OmitNil().Update()
			if err != nil {
				g.Log().Errorf(ctx, "UpdateUserDefaultSubscriptionForPaymentSuccess err:%s", err.Error())
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
