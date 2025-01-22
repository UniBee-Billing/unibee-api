package query

import (
	"context"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetRefundByRefundId(ctx context.Context, refundId string) (one *entity.Refund) {
	if len(refundId) == 0 {
		return nil
	}
	err := dao.Refund.Ctx(ctx).Where(dao.Refund.Columns().RefundId, refundId).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetRefundByGatewayRefundId(ctx context.Context, gatewayRefundId string) (one *entity.Refund) {
	if len(gatewayRefundId) == 0 {
		return nil
	}
	err := dao.Refund.Ctx(ctx).Where(dao.Refund.Columns().GatewayRefundId, gatewayRefundId).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetPendingGatewayTypeRefundsByPaymentId(ctx context.Context, paymentId string) (list []*entity.Refund) {
	if len(paymentId) == 0 {
		return nil
	}
	err := dao.Refund.Ctx(ctx).
		Where(dao.Refund.Columns().PaymentId, paymentId).
		Where(dao.Refund.Columns().Status, consts.RefundCreated).
		Where(dao.Refund.Columns().Type, consts.RefundTypeGateway).
		OmitEmpty().Scan(&list)
	if err != nil {
		list = make([]*entity.Refund, 0)
	}
	return
}
