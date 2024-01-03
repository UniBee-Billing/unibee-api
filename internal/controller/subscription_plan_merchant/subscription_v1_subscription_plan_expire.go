package subscription_plan_merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/subscription_plan_merchant/v1"
)

func (c *ControllerV1) SubscriptionPlanExpire(ctx context.Context, req *v1.SubscriptionPlanExpireReq) (res *v1.SubscriptionPlanExpireRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
