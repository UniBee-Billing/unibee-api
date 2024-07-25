package query

import (
	"context"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetUserDiscountById(ctx context.Context, id int64) (one *entity.MerchantUserDiscountCode) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantUserDiscountCode.Ctx(ctx).Where(dao.MerchantUserDiscountCode.Columns().Id, id).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
