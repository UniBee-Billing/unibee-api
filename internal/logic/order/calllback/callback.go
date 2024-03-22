package calllback

import (
	"context"
	entity "unibee/internal/model/entity/oversea_pay"
)

type MerchantOneTimePaymentCallback struct {
}

func (m MerchantOneTimePaymentCallback) PaymentRefundCancelCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	//TODO implement me
}

func (m MerchantOneTimePaymentCallback) PaymentRefundCreateCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	//TODO implement me

}

func (m MerchantOneTimePaymentCallback) PaymentRefundSuccessCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	//TODO implement me

}

func (m MerchantOneTimePaymentCallback) PaymentRefundFailureCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	//TODO implement me

}

func (m MerchantOneTimePaymentCallback) PaymentRefundReverseCallback(ctx context.Context, payment *entity.Payment, refund *entity.Refund) {
	//TODO implement me

}

func (m MerchantOneTimePaymentCallback) PaymentNeedAuthorisedCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me

}

func (m MerchantOneTimePaymentCallback) PaymentCancelCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me

}

func (m MerchantOneTimePaymentCallback) PaymentCreateCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me

}

func (m MerchantOneTimePaymentCallback) PaymentSuccessCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me

}

func (m MerchantOneTimePaymentCallback) PaymentFailureCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me

}
