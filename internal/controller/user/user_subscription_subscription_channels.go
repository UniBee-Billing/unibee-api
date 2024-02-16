package user

import (
	"context"
	"unibee-api/api/user/subscription"
	"unibee-api/internal/logic/subscription/service"
	"unibee-api/internal/query"
)

func (c *ControllerSubscription) SubscriptionChannels(ctx context.Context, req *subscription.SubscriptionChannelsReq) (res *subscription.SubscriptionChannelsRes, err error) {
	data := query.GetListSubscriptionTypeGateways(ctx)
	return &subscription.SubscriptionChannelsRes{
		Gateways: service.ConvertChannelsToRos(data),
	}, nil
}
