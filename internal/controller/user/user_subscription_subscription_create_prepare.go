package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/user/subscription"
)

func (c *ControllerSubscription) SubscriptionCreatePrepare(ctx context.Context, req *subscription.SubscriptionCreatePrepareReq) (res *subscription.SubscriptionCreatePrepareRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
