package callback

import (
	"context"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type Invalid struct {
}

func (i Invalid) PaymentNeedAuthorisedCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) PaymentCancelCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) PaymentCreateCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) PaymentSuccessCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) PaymentFailureCallback(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) {
	//TODO implement me
	panic("implement me")
}
