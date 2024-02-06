package query

import (
	"context"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/logic/gateway/ro"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

func GetGatewayPlan(ctx context.Context, planId int64, gatewayId int64) (one *entity.GatewayPlan) {
	if planId <= 0 || gatewayId <= 0 {
		return nil
	}
	err := dao.GatewayPlan.Ctx(ctx).Where(entity.GatewayPlan{PlanId: planId, GatewayId: gatewayId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetActiveGatewayPlan(ctx context.Context, planId int64, gatewayId int64) (one *entity.GatewayPlan) {
	if planId <= 0 || gatewayId <= 0 {
		return nil
	}
	err := dao.GatewayPlan.Ctx(ctx).Where(entity.GatewayPlan{PlanId: planId, GatewayId: gatewayId, Status: consts.GatewayPlanStatusActive}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetListActiveGatewayPlans(ctx context.Context, planId int64) (list []*entity.GatewayPlan) {
	if planId <= 0 {
		return nil
	}
	err := dao.GatewayPlan.Ctx(ctx).Where(entity.GatewayPlan{PlanId: planId, Status: consts.GatewayPlanStatusActive}).OmitEmpty().Scan(&list)
	if err != nil {
		list = nil
	}
	return
}

func GetListActiveOutGatewayRos(ctx context.Context, planId int64) []*ro.OutGatewayRo {
	if planId <= 0 {
		return nil
	}
	var list []*entity.GatewayPlan
	err := dao.GatewayPlan.Ctx(ctx).Where(entity.GatewayPlan{PlanId: planId, Status: consts.GatewayPlanStatusActive}).OmitEmpty().Scan(&list)
	if err != nil {
		return nil
	}
	var gateways []*ro.OutGatewayRo
	for _, one := range list {
		if one.Status == consts.GatewayPlanStatusActive {
			outChannel := GetGatewayById(ctx, one.GatewayId)
			if outChannel != nil {
				gateways = append(gateways, &ro.OutGatewayRo{
					GatewayId:   outChannel.Id,
					GatewayName: outChannel.Name,
				})
			}
		}
	}
	return gateways
}
