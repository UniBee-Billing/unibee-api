package user

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/api/user/subscription"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func (c *ControllerSubscription) SubscriptionChannels(ctx context.Context, req *subscription.SubscriptionChannelsReq) (res *subscription.SubscriptionChannelsRes, err error) {
	data := query.GetListSubscriptionTypePayChannels(ctx)
	utility.SuccessJsonExit(g.RequestFromCtx(ctx), data)
	return
}
