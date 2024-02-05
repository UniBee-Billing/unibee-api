package user

import (
	"context"
	"go-oversea-pay/api/user/subscription"
	"go-oversea-pay/internal/logic/subscription/service"
	"go-oversea-pay/internal/query"
)

func (c *ControllerSubscription) SubscriptionChannels(ctx context.Context, req *subscription.SubscriptionChannelsReq) (res *subscription.SubscriptionChannelsRes, err error) {
	data := query.GetListSubscriptionTypeGateways(ctx)
	return &subscription.SubscriptionChannelsRes{
		Gateways: service.ConvertChannelsToRos(data),
	}, nil
}
