package testclock

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/logic/subscription/billingcycle/cycle"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

type TestClockWalkRes struct {
	Walks []*cycle.BillingCycleWalkRes
	Error error
}

func WalkSubscriptionToTestClock(ctx context.Context, subId string, newTestClock int64) (*TestClockWalkRes, error) {
	if consts.GetConfigInstance().IsProd() {
		return nil, gerror.New("Test Does Not Work For Prod Env")
	}
	//TestClock Verify
	utility.Assert(len(subId) > 0, "Invalid SubscriptionId")
	utility.Assert(newTestClock > 0, "Invalid TestClock")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subId)
	utility.Assert(sub != nil, "Subscription Not Found")
	utility.Assert(sub.Status != consts.SubStatusExpired && sub.Status != consts.SubStatusCancelled, "Subscription Has Walk To End Status, Cancel or Expire")
	utility.Assert(sub.TestClock < newTestClock, "The Subscription Has Walk To The TestClock Exceed The New One")

	//firstEnd := subscription.GetPeriodEndFromStart(ctx,utility.MaxInt64(sub.CurrentPeriodEnd,sub.TrialEnd), uint64(sub.PlanId))
	// Verify Farthest Time Which Test Clock Can Set, The Max Number Of Subscription Billing Cycle Which TestClock Can Cover is 2
	var maxTimeCap int64 = 24 * 60 * 60 * 7 // Max 7days TestClock Cap
	if sub.TestClock > 0 {
		utility.Assert((newTestClock-sub.TestClock) < maxTimeCap, "TimeCap Should Lower Then 7 Days")
	}
	var result *TestClockWalkRes = &TestClockWalkRes{
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
		} else if walk.WalkHasDeal == false {
			result.Walks = append(result.Walks, walk)
			break
		} else {
			result.Walks = append(result.Walks, walk)
		}
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

	return result, nil
}
