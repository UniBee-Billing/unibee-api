package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetRefundByMerchantRefundNo(ctx context.Context, merchantRefundNo string) (one *entity.Refund) {
	err := dao.Refund.Ctx(ctx).Where(entity.Refund{OutRefundNo: merchantRefundNo}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
