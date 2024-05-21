package user

import (
	"context"
	_interface "unibee/internal/interface"
	user2 "unibee/internal/logic/user"

	"unibee/api/user/profile"
)

func (c *ControllerProfile) ChangeGateway(ctx context.Context, req *profile.ChangeGatewayReq) (res *profile.ChangeGatewayRes, err error) {
	user2.UpdateUserDefaultGatewayPaymentMethod(ctx, _interface.Context().Get(ctx).User.Id, req.GatewayId, req.PaymentMethodId)
	return &profile.ChangeGatewayRes{}, nil
}
