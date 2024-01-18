package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetRefundByRefundId(ctx context.Context, refundId string) (one *entity.Refund) {
	err := dao.Refund.Ctx(ctx).Where(entity.Refund{RefundId: refundId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
