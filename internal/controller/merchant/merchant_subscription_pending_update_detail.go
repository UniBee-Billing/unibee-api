package merchant

import (
	"context"
	"unibee/internal/logic/subscription/service"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) PendingUpdateDetail(ctx context.Context, req *subscription.PendingUpdateDetailReq) (res *subscription.PendingUpdateDetailRes, err error) {
	return &subscription.PendingUpdateDetailRes{SubscriptionPendingUpdate: service.GetSubscriptionPendingUpdateDetailByPendingUpdateId(ctx, req.SubscriptionPendingUpdateId)}, nil
}
