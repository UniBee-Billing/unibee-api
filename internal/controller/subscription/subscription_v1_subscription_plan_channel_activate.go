package subscription

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/subscription/v1"
)

func (c *ControllerV1) SubscriptionPlanChannelActivate(ctx context.Context, req *v1.SubscriptionPlanChannelActivateReq) (res *v1.SubscriptionPlanChannelActivateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
