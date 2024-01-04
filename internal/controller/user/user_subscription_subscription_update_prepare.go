package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/user/subscription"
)

func (c *ControllerSubscription) SubscriptionUpdatePrepare(ctx context.Context, req *subscription.SubscriptionUpdatePrepareReq) (res *subscription.SubscriptionUpdatePrepareRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
