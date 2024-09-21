package _interface

import (
	"context"
	entity "unibee/internal/model/entity/default"
)

type PaymentBizCallbackInterface interface {
	PaymentCreateCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice)
	PaymentSuccessCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice)
	PaymentFailureCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice)
	PaymentCancelCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice)
	PaymentNeedAuthorisedCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice)
	PaymentRefundCreateCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund)
	PaymentRefundSuccessCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund)
	PaymentRefundCancelCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund)
	PaymentRefundFailureCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund)
	PaymentRefundReverseCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund)
}
