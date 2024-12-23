package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/config"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) Config(ctx context.Context, req *subscription.ConfigReq) (res *subscription.ConfigRes, err error) {
	return &subscription.ConfigRes{Config: config.GetMerchantSubscriptionConfig(ctx, _interface.GetMerchantId(ctx))}, nil
}
