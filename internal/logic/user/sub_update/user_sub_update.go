package sub_update

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	dao "unibee/internal/dao/default"
	//"unibee/internal/logic/operation_log"
	"unibee/internal/query"
	"unibee/utility"
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
			} else {
				_, _ = redismq.Send(&redismq.Message{
					Topic:      redismq2.TopicUserAccountUpdate.Topic,
					Tag:        redismq2.TopicUserAccountUpdate.Tag,
					Body:       fmt.Sprintf("%d", user.Id),
					CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
				})
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
			} else {
				_, _ = redismq.Send(&redismq.Message{
					Topic:      redismq2.TopicUserAccountUpdate.Topic,
					Tag:        redismq2.TopicUserAccountUpdate.Tag,
					Body:       fmt.Sprintf("%d", user.Id),
					CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
				})
			}
		}
	}
}
