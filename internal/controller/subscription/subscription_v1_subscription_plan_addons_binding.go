package subscription

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/subscription/v1"
)

func (c *ControllerV1) SubscriptionPlanAddonsBinding(ctx context.Context, req *v1.SubscriptionPlanAddonsBindingReq) (res *v1.SubscriptionPlanAddonsBindingRes, err error) {
	one, err := service.SubscriptionPlanAddonsBinding(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.SubscriptionPlanAddonsBindingRes{Plan: one}, nil
}
