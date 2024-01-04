package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetOverseaRefundByMerchantRefundNo(ctx context.Context, merchantRefundNo string) (one *entity.OverseaRefund) {
	err := dao.OverseaRefund.Ctx(ctx).Where(entity.OverseaRefund{OutRefundNo: merchantRefundNo}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
