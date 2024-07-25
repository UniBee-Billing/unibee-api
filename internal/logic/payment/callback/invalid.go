package callback

import (
	"context"
	entity "unibee/internal/model/entity/default"
)

type Invalid struct {
}

func (i Invalid) PaymentRefundCancelCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {

}

func (i Invalid) PaymentRefundCreateCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {

}

func (i Invalid) PaymentRefundSuccessCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {

}

func (i Invalid) PaymentRefundFailureCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {

}

func (i Invalid) PaymentRefundReverseCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {

}

func (i Invalid) PaymentNeedAuthorisedCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {

}

func (i Invalid) PaymentCancelCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {

}

func (i Invalid) PaymentCreateCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {

}

func (i Invalid) PaymentSuccessCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {

}

func (i Invalid) PaymentFailureCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {

}
