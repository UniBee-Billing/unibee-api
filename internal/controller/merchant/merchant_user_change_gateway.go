package merchant

import (
	"context"
	"strconv"
	"unibee/internal/consts"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/util"
	user2 "unibee/internal/logic/user/sub_update"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/user"
)

func (c *ControllerUser) ChangeGateway(ctx context.Context, req *user.ChangeGatewayReq) (res *user.ChangeGatewayRes, err error) {
	if len(req.GatewayUserId) > 0 {
		targetUser := query.GetUserAccountById(ctx, req.UserId)
		utility.Assert(targetUser != nil, "user not found")
		gateway := query.GetGatewayById(ctx, req.GatewayId)
		utility.Assert(targetUser != nil, "gateway not found")
		utility.Assert(gateway.MerchantId == targetUser.MerchantId, "merchant not match:"+strconv.FormatUint(req.GatewayId, 10))
		gatewayUser := util.GetGatewayUser(ctx, req.UserId, req.GatewayId)
		if gatewayUser != nil {
			utility.Assert(gatewayUser.GatewayUserId == req.GatewayUserId, "another gateway user exist")
		}
		gatewayInfo := api.GetGatewayServiceProvider(ctx, gateway.Id).GatewayInfo(ctx)
		if gatewayInfo != nil && gatewayInfo.AutoChargeEnabled && gateway.GatewayType == consts.GatewayTypeCard {
			gatewayUserDetail, err := api.GetGatewayServiceProvider(ctx, gateway.Id).GatewayUserDetailQuery(ctx, gateway, req.GatewayUserId)
			utility.AssertError(err, "query user from gateway failed")
			utility.Assert(gatewayUserDetail != nil && gatewayUserDetail.GatewayUserId == req.GatewayUserId, "gateway user not found")
			//if len(req.PaymentMethodId) > 0 {
			//	paymentMethod := method.QueryPaymentMethod(ctx, targetUser.MerchantId, targetUser.Id, req.GatewayId, req.PaymentMethodId)
			//	utility.Assert(paymentMethod != nil, "paymentMethod not found")
			//}
			_, err = util.CreateOrUpdateGatewayUser(ctx, req.UserId, req.GatewayId, req.GatewayUserId, "")
			utility.AssertError(err, "CreateOrUpdateGatewayUser failed")
		}
	}
	user2.UpdateUserDefaultGatewayPaymentMethod(ctx, req.UserId, req.GatewayId, req.PaymentMethodId, req.GatewayPaymentType)
	return &user.ChangeGatewayRes{}, nil
}
