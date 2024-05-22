package user

import (
	"context"
	"strconv"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	user2 "unibee/internal/logic/user"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/user/profile"
)

func (c *ControllerProfile) ChangeGateway(ctx context.Context, req *profile.ChangeGatewayReq) (res *profile.ChangeGatewayRes, err error) {
	user := query.GetUserAccountById(ctx, _interface.Context().Get(ctx).User.Id)
	utility.Assert(user != nil, "user not found")
	if len(user.GatewayId) > 0 {
		oldGatewayId, err := strconv.ParseUint(user.GatewayId, 10, 64)
		if err == nil {
			gateway := query.GetGatewayById(ctx, oldGatewayId)
			newGateway := query.GetGatewayById(ctx, req.GatewayId)
			if oldGatewayId != req.GatewayId {
				utility.Assert(gateway.GatewayType != consts.GatewayTypeWireTransfer, "Can't change gateway from wire transfer to other, Please contact billing admin")
				utility.Assert(newGateway.GatewayType != consts.GatewayTypeWireTransfer, "Can't change gateway to wire transfer, Please contact billing admin")
			}
		}
	}
	user2.UpdateUserDefaultGatewayPaymentMethod(ctx, _interface.Context().Get(ctx).User.Id, req.GatewayId, req.PaymentMethodId)
	return &profile.ChangeGatewayRes{}, nil
}
