package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetPaymentById(ctx context.Context, payId int64) (one *entity.Payment) {
	err := dao.Payment.Ctx(ctx).Where(entity.Payment{Id: payId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetPaymentByPaymentId(ctx context.Context, paymentId string) (one *entity.Payment) {
	err := dao.Payment.Ctx(ctx).Where(entity.Payment{PaymentId: paymentId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetPaymentByChannelPaymentId(ctx context.Context, channelPaymentId string) (one *entity.Payment) {
	err := dao.Payment.Ctx(ctx).Where(entity.Invoice{ChannelPaymentId: channelPaymentId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetPaymentTimeLineByUniqueId(ctx context.Context, uniqueId string) (one *entity.PaymentTimeline) {
	err := dao.PaymentTimeline.Ctx(ctx).Where(entity.PaymentTimeline{UniqueId: uniqueId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
