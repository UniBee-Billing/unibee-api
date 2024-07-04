package query

import (
	"context"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func GetPlanById(ctx context.Context, id uint64) (one *entity.Plan) {
	if id <= 0 {
		return nil
	}
	err := dao.Plan.Ctx(ctx).Where(dao.Plan.Columns().Id, id).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetPlanByExternalPlanId(ctx context.Context, merchantId uint64, externalPlanId string) (one *entity.Plan) {
	if merchantId <= 0 {
		return nil
	}
	if len(externalPlanId) <= 0 {
		return nil
	}
	err := dao.Plan.Ctx(ctx).Where(dao.Plan.Columns().ExternalPlanId, externalPlanId).Where(dao.Plan.Columns().MerchantId, merchantId).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetPlansByIds(ctx context.Context, ids []int64) (list []*entity.Plan) {
	err := dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, ids).OmitEmpty().Scan(&list)
	if err != nil {
		return nil
	}
	return list
}

func GetAddonsByIds(ctx context.Context, addonIdsList []int64) (list []*entity.Plan) {
	err := dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, addonIdsList).Scan(&list)
	if err != nil {
		return nil
	}
	return list
}
