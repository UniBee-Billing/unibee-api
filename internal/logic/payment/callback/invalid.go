package callback

import (
	"context"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type Invalid struct {
}

func (i Invalid) PaymentSuccessCallback(ctx context.Context, payment *entity.Payment) {
}

func (i Invalid) PaymentFailureCallback(ctx context.Context, payment *entity.Payment) {
}
