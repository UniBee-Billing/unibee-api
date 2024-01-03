package subscription

import (
	"context"
	"go-oversea-pay/api/subscription/v1"
	"go-oversea-pay/internal/logic/subscription/service"
)

func (c *ControllerV1) SubscriptionUpdate(ctx context.Context, req *v1.SubscriptionUpdateReq) (res *v1.SubscriptionUpdateRes, err error) {
	one, err := service.SubscriptionUpdate(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.SubscriptionUpdateRes{SubscriptionPendingUpdate: one}, nil
}
