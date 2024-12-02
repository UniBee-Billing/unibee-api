package merchant

import (
	"context"
	"unibee/internal/logic/user/sub_update"

	"unibee/api/merchant/user"
)

func (c *ControllerUser) ClearAutoChargeMethod(ctx context.Context, req *user.ClearAutoChargeMethodReq) (res *user.ClearAutoChargeMethodRes, err error) {
	sub_update.ClearUserDefaultGatewayMethodForAutoCharge(ctx, req.UserId)
	return &user.ClearAutoChargeMethodRes{}, nil
}
