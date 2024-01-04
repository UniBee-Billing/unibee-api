package user

import (
	"context"
	"go-oversea-pay/api/user/subscription"
	"go-oversea-pay/internal/logic/subscription/service"
)

func (c *ControllerSubscription) SubscriptionDetail(ctx context.Context, req *subscription.SubscriptionDetailReq) (res *subscription.SubscriptionDetailRes, err error) {
	return service.SubscriptionDetail(ctx, req.SubscriptionId)
}
