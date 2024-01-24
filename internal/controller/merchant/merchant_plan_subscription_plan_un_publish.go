package merchant

import (
	"context"
	"go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/internal/logic/plan/service"
)

func (c *ControllerPlan) SubscriptionPlanUnPublish(ctx context.Context, req *plan.SubscriptionPlanUnPublishReq) (res *plan.SubscriptionPlanUnPublishRes, err error) {
	err = service.SubscriptionPlanUnPublish(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &plan.SubscriptionPlanUnPublishRes{}, nil
}
