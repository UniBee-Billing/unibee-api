package subscription

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/subscription/v1"
)

func (c *ControllerV1) SubscriptionPlanList(ctx context.Context, req *v1.SubscriptionPlanListReq) (res *v1.SubscriptionPlanListRes, err error) {
	return &v1.SubscriptionPlanListRes{Plans: service.SubscriptionPlanList(ctx, req)}, nil
}
