package service

import (
	"context"
	"unibee/internal/consts"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/payment/handler"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func PaymentGatewayCancel(ctx context.Context, payment *entity.Payment) (err error) {
	utility.Assert(payment != nil, "payment not found")
	utility.Assert(payment.Status != consts.PaymentFailed, "payment already failure")
	utility.Assert(payment.Status != consts.PaymentCancelled, "payment already cancelled")
	utility.Assert(payment.Status == consts.PaymentCreated, "payment not created status")
	utility.Assert(payment.AuthorizeStatus < consts.CaptureRequest, "payment has capture request")
	res, err := api.GetGatewayServiceProvider(ctx, payment.GatewayId).GatewayCancel(ctx, payment)
	if err != nil {
		return err
	}
	if res.Status == consts.PaymentCancelled {
		err = handler.HandlePayCancel(ctx, &handler.HandlePayReq{
			PaymentId:     payment.PaymentId,
			PayStatusEnum: consts.PaymentCancelled,
			Reason:        "Merchant Cancel",
		})
		if err != nil {
			return err
		}
	}
	return
}

func PaymentRefundGatewayCancel(ctx context.Context, refund *entity.Refund) (err error) {
	utility.Assert(refund != nil, "refund not found")
	utility.Assert(refund.Status != consts.RefundFailed, "refund already failure")
	utility.Assert(refund.Status != consts.RefundCancelled, "refund already cancelled")
	utility.Assert(refund.Status == consts.RefundCreated, "refund not created status")
	payment := query.GetPaymentByPaymentId(ctx, refund.PaymentId)
	utility.Assert(payment != nil, "payment not found")
	res, err := api.GetGatewayServiceProvider(ctx, payment.GatewayId).GatewayRefundCancel(ctx, payment, refund)
	if err != nil {
		return err
	}
	if res.Status == consts.RefundCancelled {
		err = handler.HandleRefundCancelled(ctx, &handler.HandleRefundReq{
			RefundId:         refund.RefundId,
			RefundStatusEnum: consts.RefundCancelled,
			Reason:           "Merchant Cancel",
		})
		if err != nil {
			return err
		}
	}
	return
}
