package user

import (
	"context"
	"go-oversea-pay/internal/consts"

	"go-oversea-pay/api/user/subscription"
)

func (c *ControllerSubscription) SubscriptionPayCheck(ctx context.Context, req *subscription.SubscriptionPayCheckReq) (res *subscription.SubscriptionPayCheckRes, err error) {
	detail, err := c.SubscriptionDetail(ctx, &subscription.SubscriptionDetailReq{
		SubscriptionId: req.SubscriptionId,
	})
	if err != nil {
		return nil, err
	}

	return &subscription.SubscriptionPayCheckRes{
		PayStatus:    consts.SubscriptionStatusEnum(detail.Subscription.Status),
		Subscription: detail.Subscription,
	}, nil
}
