package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"go-oversea-pay/api/merchant/plan"
)

func (c *ControllerPlan) SubscriptionPlanExpire(ctx context.Context, req *plan.SubscriptionPlanExpireReq) (res *plan.SubscriptionPlanExpireRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
