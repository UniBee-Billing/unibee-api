package merchant

import (
	"context"
	"fmt"
	v1 "go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/internal/logic/subscription/service"
)

func (c *ControllerPlan) SubscriptionPlanList(ctx context.Context, req *v1.SubscriptionPlanListReq) (res *v1.SubscriptionPlanListRes, err error) {
	fmt.Println("context: ", ctx)
	return &v1.SubscriptionPlanListRes{Plans: service.SubscriptionPlanList(ctx, req)}, nil
}
