package user

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/ro"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/user/payment"
)

func (c *ControllerPayment) MethodList(ctx context.Context, req *payment.MethodListReq) (res *payment.MethodListRes, err error) {
	merchant := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	utility.Assert(merchant != nil, "merchant not found")
	utility.Assert(req.GatewayId > 0, "invalid gatewayId")
	gateway := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(merchant.Id == gateway.MerchantId, "wrong gateway")
	gatewayUser := query.GetGatewayUser(ctx, int64(_interface.BizCtx().Get(ctx).User.Id), req.GatewayId)
	if gatewayUser != nil {
		listQuery, err := api.GetGatewayServiceProvider(ctx, req.GatewayId).GatewayUserPaymentMethodListQuery(ctx, gateway, gatewayUser.UserId)
		if err != nil {
			return nil, err
		}
		return &payment.MethodListRes{MethodList: listQuery.PaymentMethods}, nil
	} else {
		return &payment.MethodListRes{MethodList: make([]*ro.PaymentMethod, 0)}, nil
	}
}
