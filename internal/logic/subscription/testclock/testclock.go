package testclock

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee-api/internal/consts"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

type TestClockWalkRes struct {
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
	// Todo Mark Verify Farthest Time Which Test Clock Can Set, The Max Number Of Subscription Billing Cycle Which TestClock Can Cover is 2
	return &TestClockWalkRes{}, nil
}
