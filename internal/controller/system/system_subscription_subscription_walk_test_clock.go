package system

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"unibee/internal/logic/subscription/billingcycle/testclock"
	"unibee/utility"

	"unibee/api/system/subscription"
)

func (c *ControllerSubscription) SubscriptionWalkTestClock(ctx context.Context, req *subscription.SubscriptionWalkTestClockReq) (res *subscription.SubscriptionWalkTestClockRes, err error) {
	clock, err := testclock.WalkSubscriptionToTestClock(ctx, req.SubscriptionId, req.NewTestClock)
	if err != nil {
		return nil, err
	}
	glog.Infof(ctx, "SubscriptionWalkTestClock SubId:%s Res:%s", req.SubscriptionId, utility.MarshalToJsonString(clock))
	return &subscription.SubscriptionWalkTestClockRes{}, nil
}
