package _interface

import (
	"context"
	entity "unibee/internal/model/entity/oversea_pay"
)

type PaymentBizCallbackInterface interface {
	PaymentCreateCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice)
	PaymentSuccessCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice)
	PaymentFailureCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice)
	PaymentCancelCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice)
	PaymentNeedAuthorisedCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice)
}
