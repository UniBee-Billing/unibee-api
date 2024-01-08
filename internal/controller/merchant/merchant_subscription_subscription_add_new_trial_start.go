package merchant

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/merchant/subscription"
)

func (c *ControllerSubscription) SubscriptionAddNewTrialStart(ctx context.Context, req *subscription.SubscriptionAddNewTrialStartReq) (res *subscription.SubscriptionAddNewTrialStartRes, err error) {
	err = service.SubscriptionAddNewTrailEnd(ctx, req.SubscriptionId, req.TrailEnd)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionAddNewTrialStartRes{}, nil
}
