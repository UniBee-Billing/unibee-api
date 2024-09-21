package user

import (
	"context"
	"unibee/internal/consts"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) PayCheck(ctx context.Context, req *subscription.PayCheckReq) (res *subscription.PayCheckRes, err error) {
	detail, err := c.Detail(ctx, &subscription.DetailReq{
		SubscriptionId: req.SubscriptionId,
	})
	if err != nil {
		return nil, err
	}

	return &subscription.PayCheckRes{
		PayStatus:    consts.SubscriptionStatusEnum(detail.Subscription.Status),
		Subscription: detail.Subscription,
	}, nil
}
