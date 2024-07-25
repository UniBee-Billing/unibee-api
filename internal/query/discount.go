package query

import (
	"context"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetDiscountById(ctx context.Context, id uint64) (one *entity.MerchantDiscountCode) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantDiscountCode.Ctx(ctx).Where(dao.MerchantDiscountCode.Columns().Id, id).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetDiscountByCode(ctx context.Context, merchantId uint64, code string) (one *entity.MerchantDiscountCode) {
	if len(code) <= 0 || merchantId <= 0 {
		return nil
	}
	err := dao.MerchantDiscountCode.Ctx(ctx).
		Where(dao.MerchantDiscountCode.Columns().MerchantId, merchantId).
		Where(dao.MerchantDiscountCode.Columns().Code, code).
		Where(dao.MerchantDiscountCode.Columns().IsDeleted, 0).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
