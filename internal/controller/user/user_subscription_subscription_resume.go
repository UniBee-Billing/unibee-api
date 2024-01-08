package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/user/subscription"
)

func (c *ControllerSubscription) SubscriptionResume(ctx context.Context, req *subscription.SubscriptionResumeReq) (res *subscription.SubscriptionResumeRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
