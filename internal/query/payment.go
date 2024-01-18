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

func GetPaymentByMerchantOrderNo(ctx context.Context, merchantOrderNo string) (one *entity.Payment) {
	err := dao.Payment.Ctx(ctx).Where(entity.Payment{MerchantOrderNo: merchantOrderNo}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
