package merchant

import (
	"context"
	"unibee/api/merchant/subscription"
	"unibee/internal/logic/subscription/service"
)

func (c *ControllerSubscription) AddNewTrialStart(ctx context.Context, req *subscription.AddNewTrialStartReq) (res *subscription.AddNewTrialStartRes, err error) {
	err = service.SubscriptionAddNewTrialEnd(ctx, req.SubscriptionId, req.AppendTrialEndHour)
	if err != nil {
		return nil, err
	}
	return &subscription.AddNewTrialStartRes{}, nil
}
