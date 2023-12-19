package mock

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/mock/v1"
)

func (c *ControllerV1) SubscriptionPlanDetail(ctx context.Context, req *v1.SubscriptionPlanDetailReq) (res *v1.SubscriptionPlanDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
