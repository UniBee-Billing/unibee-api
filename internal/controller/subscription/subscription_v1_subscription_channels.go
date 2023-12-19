package subscription

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/subscription/v1"
)

func (c *ControllerV1) SubscriptionChannels(ctx context.Context, req *v1.SubscriptionChannelsReq) (res *v1.SubscriptionChannelsRes, err error) {
	data := query.GetListSubscriptionTypePayChannels(ctx)
	utility.SuccessJsonExit(g.RequestFromCtx(ctx), data)
	return
}
