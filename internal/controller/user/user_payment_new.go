package user

import (
	"context"
	merchantPaymentApi "unibee/api/merchant/payment"
	"unibee/api/user/payment"
	merchant "unibee/internal/controller/merchant"
	_interface "unibee/internal/interface"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerPayment) New(ctx context.Context, req *payment.NewReq) (res *payment.NewRes, err error) {
	one := query.GetUserAccountById(ctx, _interface.Context().Get(ctx).User.Id)
	utility.Assert(one != nil, "user not found")
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "plan not found")
	req.Quantity = utility.MaxInt64(1, req.Quantity)
	controllerPayment := merchant.ControllerPayment{}
	req.Metadata["PlanId"] = req.PlanId
	req.Metadata["Quantity"] = req.Quantity
	paymentRes, paymentErr := controllerPayment.New(ctx, &merchantPaymentApi.NewReq{
		UserId:      one.Id,
		Currency:    plan.Currency,
		TotalAmount: plan.Amount * req.Quantity,
		GatewayId:   req.GatewayId,
		RedirectUrl: req.RedirectUrl,
		CountryCode: one.CountryCode,
		Name:        plan.PlanName,
		Description: plan.Description,
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
