package sub

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"time"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/subscription/billingcycle/cycle"
	"unibee/internal/logic/subscription/service"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func mainTask(ctx context.Context) {
	//3 Min Invoice Out Of Pay Email
	//Subscription Cycle Email
	//Invoice 3 Day Out Of Pay Email
}

func TaskForSubscriptionBillingCycleDunningInvoice(ctx context.Context, taskName string) {
	g.Log().Debug(ctx, taskName, "Start......")
	var timeNow = gtime.Now().Timestamp()

	var subs []*entity.Subscription
	var sortKey = "task_time asc"
	var status = []int{consts.SubStatusPending, consts.SubStatusProcessing, consts.SubStatusActive, consts.SubStatusIncomplete}
	// query sub which dunningTime expired
	q := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		WhereLT(dao.Subscription.Columns().DunningTime, timeNow). //  dunning < now
		Where(dao.Subscription.Columns().Type, consts.SubTypeUniBeeControl).
		WhereIn(dao.Subscription.Columns().Status, status)
	if !config.GetConfigInstance().IsProd() {
		// Test Clock Not Enable For Prod Env
		q = q.Where(dao.Subscription.Columns().TestClock, 0)
	}
	err := q.Limit(0, 10).
		Order(sortKey).
		OmitEmpty().Scan(&subs)
	if err != nil {
		g.Log().Errorf(ctx, "%s Error:%s", taskName, err.Error())
		return
	}

	for _, sub := range subs {
		walk, err := cycle.SubPipeBillingCycleWalk(ctx, sub.SubscriptionId, timeNow, taskName)
		if err != nil {
			g.Log().Errorf(ctx, "TaskForSubscriptionBillingCycleDunningInvoice SubPipeBillingCycleWalk error:%s", err.Error())
		}
		g.Log().Infof(ctx, "TaskForSubscriptionBillingCycleDunningInvoice SubPipeBillingCycleWalk SubId:%s WalkResult:%s", sub.SubscriptionId, utility.MarshalToJsonString(walk))
		time.Sleep(2 * time.Second)
	}

	g.Log().Debug(ctx, taskName, "End......")
}

func TaskForSubscriptionTrackAfterCancelledOrExpired(ctx context.Context, taskName string) {
	g.Log().Debug(ctx, taskName, "TaskForSubscriptionTrackAfterCancelledOrExpired Start......")
	var timeNow = gtime.Now().Timestamp()

	var subs []*entity.Subscription
	var sortKey = "task_time asc"
	var status = []int{consts.SubStatusCancelled, consts.SubStatusExpired}
	// query sub which dunningTime expired
	q := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		WhereLT(dao.Subscription.Columns().DunningTime, timeNow).                //  dunning < now
		WhereGT(dao.Subscription.Columns().CurrentPeriodEnd, timeNow-(5*86400)). //  in 5 days
		Where(dao.Subscription.Columns().Type, consts.SubTypeUniBeeControl).
		WhereIn(dao.Subscription.Columns().Status, status)
	if !config.GetConfigInstance().IsProd() {
		// Test Clock Not Enable For Prod Env
		q = q.Where(dao.Subscription.Columns().TestClock, 0)
	}
	err := q.Limit(0, 10).
		Order(sortKey).
		OmitEmpty().Scan(&subs)
	if err != nil {
		g.Log().Errorf(ctx, "%s Error:%s", taskName, err.Error())
		return
	}

	for _, sub := range subs {
		walk, err := cycle.SubPipeBillingCycleWalk(ctx, sub.SubscriptionId, timeNow, taskName)
		if err != nil {
			g.Log().Errorf(ctx, "TaskForSubscriptionTrackAfterCancelledOrExpired error:%s", err.Error())
		}
		g.Log().Infof(ctx, "TaskForSubscriptionTrackAfterCancelledOrExpired WalkResult:%s", utility.MarshalToJsonString(walk))
		time.Sleep(2 * time.Second)
	}

	g.Log().Debug(ctx, taskName, "TaskForSubscriptionTrackAfterCancelledOrExpired End......")
}

func TaskForSubscriptionInitFailed(ctx context.Context, taskName string) {
	g.Log().Debug(ctx, taskName, "TaskForSubscriptionInitFailed Start......")
	var timeNow = gtime.Now().Timestamp()

	var subs []*entity.Subscription
	var status = []int{consts.SubStatusInit}
	// query sub which dunningTime expired
	q := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		WhereLT(dao.Subscription.Columns().CreateTime, timeNow-600). //  10 min
		Where(dao.Subscription.Columns().Type, consts.SubTypeUniBeeControl).
		WhereIn(dao.Subscription.Columns().Status, status)
	err := q.Limit(0, 10).
		OmitEmpty().Scan(&subs)
	if err != nil {
		g.Log().Errorf(ctx, "%s Error:%s", taskName, err.Error())
		return
	}

	for _, sub := range subs {
		err = service.SubscriptionCancel(ctx, sub.SubscriptionId, false, false, "CancelledByInitFailure")
		if err != nil {
			g.Log().Errorf(ctx, "TaskForSubscriptionInitFailed subId:%s error:%s", sub.SubscriptionId, err.Error())
		} else {
			g.Log().Infof(ctx, "TaskForSubscriptionTrackAfterCancelledOrExpired subId:%s", sub.SubscriptionId)
		}
		time.Sleep(2 * time.Second)
	}

}
