package subscription

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/subscription/v1"
)

func (c *ControllerV1) SubscriptionPlanChannelActivate(ctx context.Context, req *v1.SubscriptionPlanChannelActivateReq) (res *v1.SubscriptionPlanChannelActivateRes, err error) {
	err = service.SubscriptionPlanChannelActivate(ctx, req.PlanId, req.ChannelId)
	if err != nil {
		return nil, err
	}
	return &v1.SubscriptionPlanChannelActivateRes{}, nil
}
