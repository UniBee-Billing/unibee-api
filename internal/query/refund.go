package query

import (
	"context"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	entity "unibee/internal/model/entity/default"
)

func GetRefundsByRefundIds(ctx context.Context, refundIds []string) (list []*entity.Refund) {
	if len(refundIds) == 0 {
		return make([]*entity.Refund, 0)
	}
	err := dao.Refund.Ctx(ctx).WhereIn(dao.Refund.Columns().RefundId, refundIds).OmitEmpty().Scan(&list)
	if err != nil {
		return make([]*entity.Refund, 0)
	}
	return
}

func GetRefundByRefundId(ctx context.Context, refundId string) (one *entity.Refund) {
	if len(refundId) == 0 {
		return nil
	}
	one = _interface.GetRefundFromPreloadContext(ctx, refundId)
	if one != nil {
		return one
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
