package subscription

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/subscription/v1"
)

func (c *ControllerV1) SubscriptionPlanChannelDeactivate(ctx context.Context, req *v1.SubscriptionPlanChannelDeactivateReq) (res *v1.SubscriptionPlanChannelDeactivateRes, err error) {
	err = service.SubscriptionPlanChannelDeactivate(ctx, req.PlanId, req.ChannelId)
	if err != nil {
		return nil, err
	}
	return &v1.SubscriptionPlanChannelDeactivateRes{}, nil
}
