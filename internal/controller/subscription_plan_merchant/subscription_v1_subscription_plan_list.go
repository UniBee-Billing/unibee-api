package subscription_plan_merchant

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/subscription_plan_merchant/v1"
)

func (c *ControllerV1) SubscriptionPlanList(ctx context.Context, req *v1.SubscriptionPlanListReq) (res *v1.SubscriptionPlanListRes, err error) {
	return &v1.SubscriptionPlanListRes{Plans: service.SubscriptionPlanList(ctx, req)}, nil
}
