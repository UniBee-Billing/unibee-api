package mock

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/mock/v1"
)

func (c *ControllerV1) SubscriptionChannels(ctx context.Context, req *v1.SubscriptionChannelsReq) (res *v1.SubscriptionChannelsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
