package gateway

import (
	"context"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/gateway/ro"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

func GetGatewaySimplifyById(ctx context.Context, id int64) (one *ro.GatewaySimplify) {
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

func GetActiveGatewaySimplifyListByMerchantId(ctx context.Context, merchantId uint64) []*ro.GatewaySimplify {
	if merchantId <= 0 {
		return nil
	}
	var list []*entity.MerchantGateway
	err := dao.MerchantGateway.Ctx(ctx).Where(entity.MerchantGateway{MerchantId: merchantId, GatewayType: consts.GatewayTypeSubscription}).OmitEmpty().Scan(&list)
	if err != nil {
		return nil
	}
	var gateways []*ro.GatewaySimplify
	for _, one := range list {
		gateways = append(gateways, &ro.GatewaySimplify{
			Id:          one.Id,
			GatewayName: one.Name,
		})
	}
	return gateways
}

func GetActiveGatewaySimplifyList(ctx context.Context, planId uint64) []*ro.GatewaySimplify {
	if planId <= 0 {
		return nil
	}
	var list []*entity.GatewayPlan
	err := dao.GatewayPlan.Ctx(ctx).Where(entity.GatewayPlan{PlanId: planId, Status: consts.GatewayPlanStatusActive}).OmitEmpty().Scan(&list)
	if err != nil {
		return nil
	}
	var gateways []*ro.GatewaySimplify
	for _, one := range list {
		if one.Status == consts.GatewayPlanStatusActive {
			outChannel := query.GetGatewayById(ctx, one.GatewayId)
			if outChannel != nil {
				gateways = append(gateways, &ro.GatewaySimplify{
					Id:          outChannel.Id,
					GatewayName: outChannel.Name,
				})
			}
		}
	}
	return gateways
}
