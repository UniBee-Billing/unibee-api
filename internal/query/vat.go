package query

import (
	"context"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetVatNumberValidateHistory(ctx context.Context, merchantId uint64, vatNumber string) (res *entity.MerchantVatNumberVerifyHistory) {
	if merchantId <= 0 || len(vatNumber) == 0 {
		return nil
	}
	err := dao.MerchantVatNumberVerifyHistory.Ctx(ctx).
		Where(entity.MerchantVatNumberVerifyHistory{MerchantId: merchantId}).
		Where(entity.MerchantVatNumberVerifyHistory{VatNumber: vatNumber}).OmitEmpty().Scan(&res)
	if err != nil {
		return nil
	}
	return res
}
