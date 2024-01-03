package subscription_plan_merchant

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/subscription_plan_merchant/v1"
)

func (c *ControllerV1) SubscriptionPlanCreate(ctx context.Context, req *v1.SubscriptionPlanCreateReq) (res *v1.SubscriptionPlanCreateRes, err error) {
	one, err := service.SubscriptionPlanCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.SubscriptionPlanCreateRes{Plan: one}, nil
}
