package query

import (
	"context"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
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
