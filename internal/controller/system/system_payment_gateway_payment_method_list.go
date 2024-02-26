package system

import (
	"context"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/query"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/system/payment"
)

func (c *ControllerPayment) GatewayPaymentMethodList(ctx context.Context, req *payment.GatewayPaymentMethodListReq) (res *payment.GatewayPaymentMethodListRes, err error) {
	pay := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	if pay != nil {
		gateway := query.GetGatewayById(ctx, pay.GatewayId)
		if gateway != nil {
			methodRes, err := api.GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayUserPaymentMethodListQuery(ctx, gateway, pay.UserId)
			if err != nil {
				return nil, err
			}
			return &payment.GatewayPaymentMethodListRes{MethodList: methodRes.PaymentMethods}, nil
		}
	}
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
