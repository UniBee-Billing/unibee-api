package query

import (
	"context"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	entity "unibee/internal/model/entity/default"
)

func GetOneLatestMainPlanByMerchantId(ctx context.Context, merchantId uint64) (one *entity.Plan) {
	if merchantId <= 0 {
		return nil
	}
	err := dao.Plan.Ctx(ctx).
		Where(dao.Plan.Columns().MerchantId, merchantId).
		Where(dao.Plan.Columns().Type, consts.PlanTypeMain).
		Where(dao.Plan.Columns().IsDeleted, 0).
		Order("is_deleted desc, status asc").
		Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetPlanById(ctx context.Context, id uint64) (one *entity.Plan) {
	if id <= 0 {
		return nil
	}
	one = _interface.GetPlanFromPreloadContext(ctx, id)
	if one != nil {
		return one
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
	err := dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, ids).Scan(&list)
	if err != nil {
		return nil
	}
	return list
}

func GetPlansByProductId(ctx context.Context, merchantId uint64, productId int64) (list []*entity.Plan) {
	if productId <= 0 {
		q := dao.Plan.Ctx(ctx).
			Where(dao.Plan.Columns().MerchantId, merchantId)
		err := q.Where(q.Builder().
			WhereOrNull(dao.Plan.Columns().ProductId).
			WhereOr(dao.Plan.Columns().ProductId, 0)).
			OmitNil().Scan(&list)
		if err != nil {
			return nil
		}
		return list
	} else {
		err := dao.Plan.Ctx(ctx).
			Where(dao.Plan.Columns().MerchantId, merchantId).
			Where(dao.Plan.Columns().ProductId, productId).
			OmitEmpty().Scan(&list)
		if err != nil {
			return nil
		}
		return list
	}
}

func GetPlanIdsByProductId(ctx context.Context, merchantId uint64, productId int64) (ids []uint64) {
	list := GetPlansByProductId(ctx, merchantId, productId)
	for _, one := range list {
		ids = append(ids, one.Id)
	}
	return ids
}

func GetAddonsByIds(ctx context.Context, addonIdsList []int64) (list []*entity.Plan) {
	err := dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, addonIdsList).Scan(&list)
	if err != nil {
		return nil
	}
	return list
}
