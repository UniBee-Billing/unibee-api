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

func GetAddonsByIds(ctx context.Context, addonIdsList []int64) (list []*entity.Plan) {
	err := dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, addonIdsList).Scan(&list)
	if err != nil {
		return nil
	}
	return list
}
