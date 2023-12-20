package subscription

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/subscription/v1"
)

func (c *ControllerV1) SubscriptionDetail(ctx context.Context, req *v1.SubscriptionDetailReq) (res *v1.SubscriptionDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
