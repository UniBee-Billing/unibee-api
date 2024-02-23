package query

import (
	"context"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

func GetGatewayPlan(ctx context.Context, planId uint64, gatewayId int64) (one *entity.GatewayPlan) {
	if planId <= 0 || gatewayId <= 0 {
		return nil
	}
	err := dao.GatewayPlan.Ctx(ctx).Where(entity.GatewayPlan{PlanId: planId, GatewayId: gatewayId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetActiveGatewayPlan(ctx context.Context, planId uint64, gatewayId int64) (one *entity.GatewayPlan) {
	if planId <= 0 || gatewayId <= 0 {
		return nil
	}
	err := dao.GatewayPlan.Ctx(ctx).Where(entity.GatewayPlan{PlanId: planId, GatewayId: gatewayId, Status: consts.GatewayPlanStatusActive}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetListActiveGatewayPlans(ctx context.Context, planId uint64) (list []*entity.GatewayPlan) {
	if planId <= 0 {
		return nil
	}
	err := dao.GatewayPlan.Ctx(ctx).Where(entity.GatewayPlan{PlanId: planId, Status: consts.GatewayPlanStatusActive}).OmitEmpty().Scan(&list)
	if err != nil {
		list = nil
	}
	return
}
