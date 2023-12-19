package subscription

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/subscription/v1"
)

func (c *ControllerV1) SubscriptionCreate(ctx context.Context, req *v1.SubscriptionCreateReq) (res *v1.SubscriptionCreateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
