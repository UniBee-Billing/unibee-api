package subscription_plan_merchant

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/subscription_plan_merchant/v1"
)

func (c *ControllerV1) SubscriptionPlanEdit(ctx context.Context, req *v1.SubscriptionPlanEditReq) (res *v1.SubscriptionPlanEditRes, err error) {
	one, err := service.SubscriptionPlanEdit(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.SubscriptionPlanEditRes{Plan: one}, nil
}
