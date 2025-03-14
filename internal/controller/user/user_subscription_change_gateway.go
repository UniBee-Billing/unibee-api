package user

import (
	"context"
	"unibee/api/user/subscription"
	"unibee/internal/logic/subscription/handler"
)

func (c *ControllerSubscription) ChangeGateway(ctx context.Context, req *subscription.ChangeGatewayReq) (res *subscription.ChangeGatewayRes, err error) {
	_, err = handler.ChangeSubscriptionGateway(ctx, req.SubscriptionId, req.GatewayId, "", req.PaymentMethodId)
	if err != nil {
		return nil, err
	}
	return &subscription.ChangeGatewayRes{}, nil
}
