package _interface

import (
	"context"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type PaymentBizCallbackInterface interface {
	PaymentSuccessCallback(ctx context.Context, payment *entity.Payment)
	PaymentFailureCallback(ctx context.Context, payment *entity.Payment)
}
