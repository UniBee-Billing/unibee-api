package testclock

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"time"
	"unibee/internal/cmd/config"
	"unibee/internal/cronjob/invoice"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/subscription/billingcycle/cycle"
	"unibee/internal/query"
	"unibee/utility"
)

type TestClockWalkRes struct {
	Walks []*cycle.BillingCycleWalkRes
	Error error
}

func WalkSubscriptionToTestClock(ctx context.Context, subId string, newTestClock int64) (*TestClockWalkRes, error) {
	//TestClock Verify
	utility.Assert(len(subId) > 0, "Invalid SubscriptionId")
	utility.Assert(newTestClock > 0, "Invalid TestClock")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subId)
	utility.Assert(sub != nil, "Subscription Not Found")
	if config.GetConfigInstance().IsProd() && sub.TestClock <= 0 {
		return nil, gerror.New("AdvanceTime Does Not Work For Prod Env")
	}
	//utility.Assert(sub.Status != consts.SubStatusExpired && sub.Status != consts.SubStatusCancelled, "Subscription Has Cancel or Expire")
	utility.Assert(sub.TestClock < newTestClock, "The Subscription Has Walk To The TestClock Exceed The New One")
	//firstEnd := subscription.GetPeriodEndFromStart(ctx,utility.MaxInt64(sub.CurrentPeriodEnd,sub.TrialEnd), uint64(sub.PlanId))
	// Verify Farthest Time Which Test Clock Can Set, The Max Number Of Subscription Billing Cycle Which TestClock Can Cover is 2
	//var maxTimeCap int64 = 60 * 24 * 60 * 60 // Max 7days TestClock Cap
	clockTimeNow := utility.MaxInt64(gtime.Timestamp(), sub.TestClock)
	//if sub.TestClock > 0 {
	plan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(plan != nil, "Plan not found")
	if strings.Compare(strings.ToLower(plan.IntervalUnit), "year") == 0 {
		utility.Assert((newTestClock-clockTimeNow) < 2*365*24*60*60, "TimeCap Should Lower Then 2 Years")
	} else {
		utility.Assert((newTestClock-clockTimeNow) < 60*24*60*60, "TimeCap Should Lower Then 60 Days")
	}
	//}
	var result = &TestClockWalkRes{
		Walks: make([]*cycle.BillingCycleWalkRes, 0),
		Error: nil,
	}
	for {
		walk, err := cycle.SubPipeBillingCycleWalk(ctx, subId, newTestClock, "WalkSubscriptionToTestClock")
		if err != nil {
			result.Error = err
			break
		} else if walk == nil {
			result.Error = gerror.Newf("walk is nil")
			break
		} else if walk.WalkUnfinished == false {
			result.Walks = append(result.Walks, walk)
			break
		} else {
			result.Walks = append(result.Walks, walk)
		}
		time.Sleep(500)
	}
	if result.Error == nil {
		_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().TestClock: newTestClock,
		}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
		if err != nil {
			g.Log().Print(ctx, "WalkSubscriptionToTestClock Update TestClock err:", err.Error())
			return nil, err
		}
	}
	g.Log().Print(ctx, "WalkSubscriptionToTestClock result:%s", utility.MarshalToJsonString(result))
	invoice.ExpireUserSubInvoices(ctx, sub, newTestClock)
	return result, nil
}
