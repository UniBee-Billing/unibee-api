package user

import (
	"context"
	"unibee/internal/logic/subscription/service"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) ChangeGateway(ctx context.Context, req *subscription.ChangeGatewayReq) (res *subscription.ChangeGatewayRes, err error) {
	err = service.ChangeSubscriptionGateway(ctx, req.SubscriptionId, req.GatewayId)
	if err != nil {
		return nil, err
	}
	return &subscription.ChangeGatewayRes{}, nil
}
