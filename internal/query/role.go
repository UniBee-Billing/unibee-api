package query

import (
	"context"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func GetRoleById(ctx context.Context, id uint64) (one *entity.MerchantRole) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantRole.Ctx(ctx).Where(dao.MerchantRole.Columns().Id, id).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetRoleByName(ctx context.Context, merchantId uint64, role string) (one *entity.MerchantRole) {
	if len(role) <= 0 || merchantId <= 0 {
		return nil
	}
	err := dao.MerchantRole.Ctx(ctx).
		Where(dao.MerchantRole.Columns().MerchantId, merchantId).
		Where(dao.MerchantRole.Columns().Role, role).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}