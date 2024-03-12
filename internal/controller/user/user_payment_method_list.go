package user

import (
	"context"
	"unibee/api/bean"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/gateway_bean"
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
		var gatewayPaymentId string
		if len(req.PaymentId) > 0 {
			one := query.GetPaymentByPaymentId(ctx, req.PaymentId)
			if one != nil {
				gatewayPaymentId = one.GatewayPaymentId
			}
		}
		listQuery, err := api.GetGatewayServiceProvider(ctx, req.GatewayId).GatewayUserPaymentMethodListQuery(ctx, gateway, &gateway_bean.GatewayUserPaymentMethodReq{
			UserId:           gatewayUser.UserId,
			GatewayPaymentId: gatewayPaymentId,
		})
		if err != nil {
			return nil, err
		}
		return &payment.MethodListRes{MethodList: listQuery.PaymentMethods}, nil
	} else {
		return &payment.MethodListRes{MethodList: make([]*bean.PaymentMethod, 0)}, nil
	}
}
