package user

import (
	"context"
	"unibee/api/user/subscription"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
)

func (c *ControllerSubscription) SubscriptionChannels(ctx context.Context, req *subscription.SubscriptionChannelsReq) (res *subscription.SubscriptionChannelsRes, err error) {
	data := query.GetListSubscriptionTypeGateways(ctx)
	return &subscription.SubscriptionChannelsRes{
		Gateways: service.ConvertChannelsToRos(data),
	}, nil
}
