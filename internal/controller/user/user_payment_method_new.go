package user

import (
	"context"
	"unibee/api/user/payment"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/payment/method"
)

func (c *ControllerPayment) MethodNew(ctx context.Context, req *payment.MethodNewReq) (res *payment.MethodNewRes, err error) {
	url, one := method.NewPaymentMethod(ctx, &method.NewPaymentMethodInternalReq{
		MerchantId:     _interface.GetMerchantId(ctx),
		UserId:         _interface.Context().Get(ctx).User.Id,
		Currency:       req.Currency,
		GatewayId:      req.GatewayId,
		SubscriptionId: req.SubscriptionId,
		RedirectUrl:    req.RedirectUrl,
		Type:           req.Type,
		Metadata:       req.Metadata,
	})
	return &payment.MethodNewRes{
		Url:    url,
		Method: one,
	}, nil
}
