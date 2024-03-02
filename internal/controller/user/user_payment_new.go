package user

import (
	"context"
	"strings"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/user/payment"
)

func (c *ControllerPayment) New(ctx context.Context, req *payment.NewReq) (res *payment.NewRes, err error) {
	merchant := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	utility.Assert(merchant != nil, "merchant not found")
	utility.Assert(req.GatewayId > 0, "invalid gatewayId")
	gateway := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(merchant.Id == gateway.MerchantId, "wrong gateway")
	utility.Assert(len(req.Type) > 0 && strings.Compare(req.Type, "card") == 0, "invalid type, should be card")
	createResult, err := api.GetGatewayServiceProvider(ctx, req.GatewayId).GatewayUserCreateAndBindPaymentMethod(ctx, gateway, int64(_interface.BizCtx().Get(ctx).User.Id), req.Data)
	if err != nil {
		return nil, err
	}
	return &payment.NewRes{Method: createResult.PaymentMethod}, nil
}
