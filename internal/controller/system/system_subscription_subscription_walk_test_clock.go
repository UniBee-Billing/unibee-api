package system

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee-api/api/system/subscription"
)

func (c *ControllerSubscription) SubscriptionWalkTestClock(ctx context.Context, req *subscription.SubscriptionWalkTestClockReq) (res *subscription.SubscriptionWalkTestClockRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
