package calllback

import (
	"context"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type MerchantOneTimePaymentCallback struct {
}

func (m MerchantOneTimePaymentCallback) PaymentNeedAuthorisedCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me
	panic("implement me")
}

func (m MerchantOneTimePaymentCallback) PaymentCancelCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me
	panic("implement me")
}

func (m MerchantOneTimePaymentCallback) PaymentCreateCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me
	panic("implement me")
}

func (m MerchantOneTimePaymentCallback) PaymentSuccessCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me
	panic("implement me")
}

func (m MerchantOneTimePaymentCallback) PaymentFailureCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me
	panic("implement me")
}
