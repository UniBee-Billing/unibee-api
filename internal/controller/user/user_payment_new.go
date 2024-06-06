package user

import (
	"context"
	merchantPaymentApi "unibee/api/merchant/payment"
	"unibee/api/user/payment"
	merchant "unibee/internal/controller/merchant"
	_interface "unibee/internal/interface"
)

func (c *ControllerPayment) New(ctx context.Context, req *payment.NewReq) (res *payment.NewRes, err error) {
	controllerPayment := merchant.ControllerPayment{}
	paymentRes, paymentErr := controllerPayment.New(ctx, &merchantPaymentApi.NewReq{
		UserId:      _interface.Context().Get(ctx).User.Id,
		Currency:    req.Currency,
		TotalAmount: req.TotalAmount,
		GatewayId:   req.GatewayId,
		RedirectUrl: req.RedirectUrl,
		CountryCode: req.CountryCode,
		Name:        req.Name,
		Description: req.Description,
		Items:       req.Items,
		Metadata:    req.Metadata,
	})

	if paymentErr != nil {
		return nil, paymentErr
	}
	return &payment.NewRes{
		Status:            paymentRes.Status,
		PaymentId:         paymentRes.PaymentId,
		ExternalPaymentId: paymentRes.ExternalPaymentId,
		Link:              paymentRes.Link,
		Action:            paymentRes.Action,
	}, nil
}
