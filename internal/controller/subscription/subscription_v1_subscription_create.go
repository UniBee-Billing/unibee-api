package subscription

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/subscription/v1"
)

func (c *ControllerV1) SubscriptionCreate(ctx context.Context, req *v1.SubscriptionCreateReq) (res *v1.SubscriptionCreateRes, err error) {
	one, err := service.SubscriptionCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.SubscriptionCreateRes{Subscription: one}, nil
}
