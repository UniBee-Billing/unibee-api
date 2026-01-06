package query

import (
	"context"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetVatNumberValidateHistoryById(ctx context.Context, id int64) (res *entity.MerchantVatNumberVerifyHistory) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantVatNumberVerifyHistory.Ctx(ctx).
		Where(dao.MerchantVatNumberVerifyHistory.Columns().Id, id).
		OmitEmpty().Scan(&res)
	if err != nil {
		return nil
	}
	return res
}

func GetVatNumberValidateHistory(ctx context.Context, merchantId uint64, vatNumber string) (res *entity.MerchantVatNumberVerifyHistory) {
	if merchantId <= 0 || len(vatNumber) == 0 {
		return nil
	}
	err := dao.MerchantVatNumberVerifyHistory.Ctx(ctx).
		Where(dao.MerchantVatNumberVerifyHistory.Columns().MerchantId, merchantId).
		Where(dao.MerchantVatNumberVerifyHistory.Columns().VatNumber, vatNumber).
		OmitEmpty().Scan(&res)
	if err != nil {
		return nil
	}
	return res
}
