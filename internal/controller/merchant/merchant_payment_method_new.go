package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/payment/method"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) MethodNew(ctx context.Context, req *payment.MethodNewReq) (res *payment.MethodNewRes, err error) {
	url, one := method.NewPaymentMethod(ctx, &method.NewPaymentMethodInternalReq{
		MerchantId:     _interface.GetMerchantId(ctx),
		UserId:         req.UserId,
		Currency:       req.Currency,
		GatewayId:      req.GatewayId,
		SubscriptionId: req.SubscriptionId,
		RedirectUrl:    req.RedirectUrl,
		Type:           req.Type,
		Data:           req.Data,
	})
	return &payment.MethodNewRes{Method: one, Url: url}, nil
}
