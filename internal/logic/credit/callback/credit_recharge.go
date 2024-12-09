package callback

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/credit/recharge/handler"
	entity "unibee/internal/model/entity/default"
)

type CreditRechargeCallback struct{}

func (c CreditRechargeCallback) PaymentCreateCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	return
}

func (c CreditRechargeCallback) PaymentSuccessCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	err := handler.HandleCreditRechargeSuccess(ctx, invoice.UniqueId, invoice, payment)
	if err != nil {
		g.Log().Errorf(ctx, "PaymentSuccessCallback Credit Recharge Error:%s\n", err.Error())
		return
	}
}

func (c CreditRechargeCallback) PaymentFailureCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	err := handler.HandleCreditRechargeFailed(ctx, invoice.UniqueId)
	if err != nil {
		g.Log().Errorf(ctx, "PaymentFailureCallback Credit Recharge Error:%s\n", err.Error())
		return
	}
}

func (c CreditRechargeCallback) PaymentCancelCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	err := handler.HandleCreditRechargeFailed(ctx, invoice.UniqueId)
	if err != nil {
		g.Log().Errorf(ctx, "PaymentCancelCallback Credit Recharge Error:%s\n", err.Error())
		return
	}
}

func (c CreditRechargeCallback) PaymentNeedAuthorisedCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	return
}

func (c CreditRechargeCallback) PaymentRefundCreateCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	return
}

func (c CreditRechargeCallback) PaymentRefundSuccessCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	return
}

func (c CreditRechargeCallback) PaymentRefundCancelCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	return
}

func (c CreditRechargeCallback) PaymentRefundFailureCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	return
}

func (c CreditRechargeCallback) PaymentRefundReverseCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	return
}
