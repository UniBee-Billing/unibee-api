package util

import (
	"context"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetPaymentByPaymentId(ctx context.Context, paymentId string) (one *entity.Payment) {
	if len(paymentId) == 0 {
		return nil
	}
	err := dao.Payment.Ctx(ctx).Where(dao.Payment.Columns().PaymentId, paymentId).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
