package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/ro"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) MethodList(ctx context.Context, req *payment.MethodListReq) (res *payment.MethodListRes, err error) {
	merchant := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	utility.Assert(merchant != nil, "merchant not found")
	utility.Assert(req.GatewayId > 0, "invalid gatewayId")
	gateway := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(merchant.Id == gateway.MerchantId, "wrong gateway")

	var gatewayPaymentId string
	if len(req.PaymentId) > 0 {
		one := query.GetPaymentByPaymentId(ctx, req.PaymentId)
		if one != nil {
			gatewayPaymentId = one.GatewayPaymentId
		}
	}
	listQuery, err := api.GetGatewayServiceProvider(ctx, req.GatewayId).GatewayUserPaymentMethodListQuery(ctx, gateway, &ro.GatewayUserPaymentMethodReq{
		UserId:           int64(req.UserId),
		GatewayPaymentId: gatewayPaymentId,
	})
	if err != nil {
		return nil, err
	}
	return &payment.MethodListRes{MethodList: listQuery.PaymentMethods}, nil

}
