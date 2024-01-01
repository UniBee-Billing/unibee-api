package subscription

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/subscription/v1"
)

func (c *ControllerV1) SubscriptionPlanEdit(ctx context.Context, req *v1.SubscriptionPlanEditReq) (res *v1.SubscriptionPlanEditRes, err error) {
	one, err := service.SubscriptionPlanEdit(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.SubscriptionPlanEditRes{Plan: one}, nil
}
