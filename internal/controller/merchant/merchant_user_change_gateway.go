package merchant

import (
	"context"
	user2 "unibee/internal/logic/user/sub_update"

	"unibee/api/merchant/user"
)

func (c *ControllerUser) ChangeGateway(ctx context.Context, req *user.ChangeGatewayReq) (res *user.ChangeGatewayRes, err error) {
	user2.UpdateUserDefaultGatewayPaymentMethod(ctx, req.UserId, req.GatewayId, req.PaymentMethodId)
	return &user.ChangeGatewayRes{}, nil
}
