package sub_update

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/operation_log"
	"unibee/internal/query"
)

func UpdateUserDefaultSubscriptionForUpdate(ctx context.Context, userId uint64, subscriptionId string) {
	if userId > 0 && len(subscriptionId) > 0 {
		one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
		user := query.GetUserAccountById(ctx, userId)
		var subName = ""
		if one != nil && user != nil && user.SubscriptionId == subscriptionId {
			plan := query.GetPlanById(ctx, one.PlanId)

			if plan != nil {
				subName = plan.PlanName
			}
			_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().PlanId:             one.PlanId,
				dao.UserAccount.Columns().SubscriptionId:     subscriptionId,
				dao.UserAccount.Columns().SubscriptionStatus: one.Status,
				dao.UserAccount.Columns().SubscriptionName:   subName,
				dao.UserAccount.Columns().BillingType:        1,
				dao.UserAccount.Columns().GmtModify:          gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, userId).OmitNil().Update()
			if err != nil {
				g.Log().Errorf(ctx, "UpdateUserDefaultSubscriptionForUpdate err:%s", err.Error())
			}
		}
	}
}

func UpdateUserDefaultSubscriptionForPaymentSuccess(ctx context.Context, userId uint64, subscriptionId string) {
	if userId > 0 && len(subscriptionId) > 0 {
		one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
		user := query.GetUserAccountById(ctx, userId)
		var subName = ""
		if one != nil && user != nil {
			plan := query.GetPlanById(ctx, one.PlanId)
			if plan != nil {
				subName = plan.PlanName
			}
			_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().PlanId:             one.PlanId,
				dao.UserAccount.Columns().SubscriptionId:     subscriptionId,
				dao.UserAccount.Columns().SubscriptionStatus: one.Status,
				dao.UserAccount.Columns().SubscriptionName:   subName,
				dao.UserAccount.Columns().BillingType:        1,
				dao.UserAccount.Columns().GmtModify:          gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, userId).OmitNil().Update()
			if err != nil {
				g.Log().Errorf(ctx, "UpdateUserDefaultSubscriptionForPaymentSuccess err:%s", err.Error())
			}
		}
	}
}

func UpdateUserDefaultVatNumber(ctx context.Context, userId uint64, vatNumber string) {
	if userId > 0 && len(vatNumber) > 0 {
		user := query.GetUserAccountById(ctx, userId)
		if user == nil {
			return
		}
		_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
			dao.UserAccount.Columns().VATNumber: vatNumber,
			dao.UserAccount.Columns().GmtModify: gtime.Now(),
		}).Where(dao.UserAccount.Columns().Id, userId).OmitNil().Update()
		if err != nil {
			g.Log().Errorf(ctx, "UpdateUserDefaultVatNumber err:%s", err.Error())
		}

		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     user.MerchantId,
			Target:         fmt.Sprintf("User(%v)", user.Id),
			Content:        fmt.Sprintf("UpdateVATNumber(%s)", vatNumber),
			UserId:         user.Id,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, nil)
	}
}
