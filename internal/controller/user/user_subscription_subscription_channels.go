package user

import (
	"context"
	"go-oversea-pay/api/user/subscription"
	"go-oversea-pay/internal/logic/subscription/service"
	"go-oversea-pay/internal/query"
)

func (c *ControllerSubscription) SubscriptionChannels(ctx context.Context, req *subscription.SubscriptionChannelsReq) (res *subscription.SubscriptionChannelsRes, err error) {
	data := query.GetListSubscriptionTypePayChannels(ctx)
	return &subscription.SubscriptionChannelsRes{
		Channels: service.ConvertChannelsToRos(data),
	}, nil
}
