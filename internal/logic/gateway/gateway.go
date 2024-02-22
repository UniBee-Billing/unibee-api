package gateway

import (
	"context"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/logic/gateway/ro"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
)

func GetOutGatewayRoById(ctx context.Context, id int64) (one *ro.OutGatewayRo) {
	if id <= 0 {
		return nil
	}
	m := dao.MerchantGateway.Ctx(ctx)
	err := m.Where(entity.MerchantGateway{Id: uint64(id)}).OmitEmpty().Scan(&one)
	if err != nil || one == nil {
		return nil
	}
	return one
}

func GetListActiveOutGatewayRosByMerchantId(ctx context.Context, merchantId uint64) []*ro.OutGatewayRo {
	if merchantId <= 0 {
		return nil
	}
	var list []*entity.MerchantGateway
	err := dao.MerchantGateway.Ctx(ctx).Where(entity.MerchantGateway{MerchantId: merchantId, GatewayType: consts.GatewayTypeSubscription}).OmitEmpty().Scan(&list)
	if err != nil {
		return nil
	}
	var gateways []*ro.OutGatewayRo
	for _, one := range list {
		gateways = append(gateways, &ro.OutGatewayRo{
			Id:          one.Id,
			GatewayName: one.Name,
		})
	}
	return gateways
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
			outChannel := query.GetGatewayById(ctx, one.GatewayId)
			if outChannel != nil {
				gateways = append(gateways, &ro.OutGatewayRo{
					Id:          outChannel.Id,
					GatewayName: outChannel.Name,
				})
			}
		}
	}
	return gateways
}
