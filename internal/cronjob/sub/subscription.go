package sub

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"time"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/logic/subscription/billingcycle/cycle"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/utility"
)

func mainTask(ctx context.Context) {
	//3 Min Invoice Out Of Pay Email
	//Subscription Cycle Email
	//Invoice 3 Day Out Of Pay Email
}

func SubscriptionBillingCycleDunningInvoice(ctx context.Context, taskName string) {
	g.Log().Debug(ctx, taskName, "Start......")
	var timeNow = gtime.Now().Timestamp()

	var subs []*entity.Subscription
	var sortKey = "task_time asc"
	var status = []int{consts.SubStatusCreate, consts.SubStatusActive, consts.SubStatusIncomplete}
	// query sub which dunningTime expired
	q := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		WhereLT(dao.Subscription.Columns().DunningTime, timeNow). //  dunning < now
		Where(dao.Subscription.Columns().Type, consts.SubTypeUniBeeControl).
		WhereIn(dao.Subscription.Columns().Status, status)
	if !consts.GetConfigInstance().IsProd() {
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
			g.Log().Errorf(ctx, "SubscriptionBillingCycleDunningInvoice SubPipeBillingCycleWalk error:%s", err.Error())
		}
		g.Log().Errorf(ctx, "SubscriptionBillingCycleDunningInvoice SubPipeBillingCycleWalk WalkResult:%s", utility.MarshalToJsonString(walk))
		time.Sleep(10 * time.Second)
	}

	g.Log().Debug(ctx, taskName, "End......")
}
