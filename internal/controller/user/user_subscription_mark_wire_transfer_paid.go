package user

import (
	"context"
	"unibee/internal/logic/subscription/service"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) MarkWireTransferPaid(ctx context.Context, req *subscription.MarkWireTransferPaidReq) (res *subscription.MarkWireTransferPaidRes, err error) {
	err = service.MarkSubscriptionProcessed(ctx, req.SubscriptionId)
	if err != nil {
		return nil, err
	}
	return &subscription.MarkWireTransferPaidRes{}, nil
}
