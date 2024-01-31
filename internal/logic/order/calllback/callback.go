package calllback

import (
	"context"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type MerchantOneTimePaymentCallback struct {
}

func (m MerchantOneTimePaymentCallback) PaymentSuccessCallback(ctx context.Context, payment *entity.Payment) {
	//TODO implement me
	panic("implement me")
}

func (m MerchantOneTimePaymentCallback) PaymentFailureCallback(ctx context.Context, payment *entity.Payment) {
	//TODO implement me
	panic("implement me")
}
