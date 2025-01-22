package util

import (
	"context"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetOtherPendingGatewayTypeRefundsByPaymentId(ctx context.Context, paymentId string, exceptRefundId string) (list []*entity.Refund) {
	if len(paymentId) == 0 {
		return nil
	}
	err := dao.Refund.Ctx(ctx).
		Where(dao.Refund.Columns().PaymentId, paymentId).
		WhereNot(dao.Refund.Columns().RefundId, exceptRefundId).
		Where(dao.Refund.Columns().Status, consts.RefundCreated).
		Where(dao.Refund.Columns().Type, consts.RefundTypeGateway).
		OmitEmpty().Scan(&list)
	if err != nil {
		list = make([]*entity.Refund, 0)
	}
	return
}
