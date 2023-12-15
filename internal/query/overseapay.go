package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetOverseaPayById(ctx context.Context, payId int64) (one *entity.OverseaPay) {
	err := dao.OverseaPay.Ctx(ctx).Where(entity.OverseaPay{Id: payId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetOverseaPayByMerchantOrderNo(ctx context.Context, merchantOrderNo string) (one *entity.OverseaPay) {
	err := dao.OverseaPay.Ctx(ctx).Where(entity.OverseaPay{MerchantOrderNo: merchantOrderNo}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
