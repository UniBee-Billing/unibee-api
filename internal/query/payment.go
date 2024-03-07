package query

import (
	"context"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/gateway/ro"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func GetPaymentById(ctx context.Context, id int64) (one *entity.Payment) {
	if id <= 0 {
		return nil
	}
	err := dao.Payment.Ctx(ctx).Where(entity.Payment{Id: id}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

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

func GetPaymentByGatewayPaymentId(ctx context.Context, gatewayPaymentId string) (one *entity.Payment) {
	if len(gatewayPaymentId) == 0 {
		return nil
	}
	err := dao.Payment.Ctx(ctx).Where(dao.Payment.Columns().GatewayPaymentId, gatewayPaymentId).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetPaymentTimeLineByUniqueId(ctx context.Context, uniqueId string) (one *entity.PaymentTimeline) {
	if len(uniqueId) == 0 {
		return nil
	}
	err := dao.PaymentTimeline.Ctx(ctx).Where(dao.PaymentTimeline.Columns().UniqueId, uniqueId).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetPaymentDetail(ctx context.Context, merchantId uint64, paymentId string) *ro.PaymentDetailRo {
	one := GetPaymentByPaymentId(ctx, paymentId)
	utility.Assert(merchantId == one.MerchantId, "merchant not match")
	if one != nil {
		return &ro.PaymentDetailRo{
			User:    ro.SimplifyUserAccount(GetUserAccountById(ctx, uint64(one.UserId))),
			Payment: ro.SimplifyPayment(one),
		}
	}
	return nil
}
