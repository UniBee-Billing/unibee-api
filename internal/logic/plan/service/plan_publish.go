package service

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/logic/gateway/api"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func SubscriptionPlanActivate(ctx context.Context, planId int64) error {
	//发布 Plan
	utility.Assert(planId > 0, "invalid planId")
	one := query.GetPlanById(ctx, planId)
	utility.Assert(one != nil, "plan not found, invalid planId")
	if one.Status == consts.PlanStatusActive {
		//已成功
		return nil
	}
	update, err := dao.SubscriptionPlan.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPlan.Columns().Status:    consts.PlanStatusActive,
		dao.SubscriptionPlan.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionPlan.Columns().Id, planId).OmitNil().Update()
	if err != nil {
		return err
	}
	affected, err := update.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return gerror.New("internal err, publish count != 1")
	}
	return nil
}

func SubscriptionPlanChannelTransferAndActivate(ctx context.Context, planId int64, gatewayId int64) error {
	intervals := []string{"day", "month", "year", "week"}
	plan := query.GetPlanById(ctx, planId)
	utility.Assert(plan != nil, "plan not found")
	utility.Assert(utility.StringContainsElement(intervals, strings.ToLower(plan.IntervalUnit)), "IntervalUnit Error，Must One Of day｜month｜year｜week")
	if strings.ToLower(plan.IntervalUnit) == "day" {
		utility.Assert(plan.IntervalCount <= 365, "IntervalCount Must Lower Then 365 While IntervalUnit is day")
	} else if strings.ToLower(plan.IntervalUnit) == "month" {
		utility.Assert(plan.IntervalCount <= 12, "IntervalCount Must Lower Then 12 While IntervalUnit is month")
	} else if strings.ToLower(plan.IntervalUnit) == "year" {
		utility.Assert(plan.IntervalCount <= 1, "IntervalCount Must Lower Then 52 While IntervalUnit is year")
	} else if strings.ToLower(plan.IntervalUnit) == "week" {
		utility.Assert(plan.IntervalCount <= 52, "IntervalCount Must Lower Then 52 While IntervalUnit is week")
	}
	gateway := query.GetSubscriptionTypeGatewayById(ctx, gatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	gatewayPlan := query.GetGatewayPlan(ctx, planId, gatewayId)
	if gatewayPlan == nil {
		gatewayPlan = &entity.GatewayPlan{
			PlanId:    planId,
			GatewayId: gatewayId,
			Status:    consts.GatewayPlanStatusInit,
			CreateAt:  gtime.Now().Timestamp(),
		}
		//保存gatewayPlan
		result, err := dao.GatewayPlan.Ctx(ctx).Data(gatewayPlan).OmitNil().Insert(gatewayPlan)
		if err != nil {
			err = gerror.Newf(`SubscriptionGatewayPlanTransferAndActivate record insert failure %s`, err)
			gatewayPlan = nil
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			gatewayPlan = nil
			return err
		}
		gatewayPlan.Id = uint64(uint(id))
	}
	if len(gatewayPlan.GatewayProductId) == 0 {
		//产品尚未创建
		if len(plan.GatewayProductName) == 0 {
			plan.GatewayProductName = plan.PlanName
		}
		if len(plan.GatewayProductDescription) == 0 {
			plan.GatewayProductDescription = plan.Description
		}
		res, err := api.GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayProductCreate(ctx, plan, gatewayPlan)
		if err != nil {
			return err
		}
		//更新 gatewayPlan
		_, err = dao.GatewayPlan.Ctx(ctx).Data(g.Map{
			dao.GatewayPlan.Columns().GatewayProductId:     res.GatewayProductId,
			dao.GatewayPlan.Columns().GatewayProductStatus: res.GatewayProductStatus,
		}).Where(dao.GatewayPlan.Columns().Id, gatewayPlan.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	return gerror.Newf("SubscriptionPlanChannelTransferAndActivate update err:%s", update)
		//}
		gatewayPlan.GatewayProductId = res.GatewayProductId
		gatewayPlan.GatewayProductStatus = res.GatewayProductStatus
	}
	if len(gatewayPlan.GatewayPlanId) == 0 {
		//创建 并激活 Plan
		res, err := api.GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayPlanCreateAndActivate(ctx, plan, gatewayPlan)
		if err != nil {
			return err
		}
		_, err = dao.GatewayPlan.Ctx(ctx).Data(g.Map{
			dao.GatewayPlan.Columns().GatewayPlanId:        res.GatewayPlanId,
			dao.GatewayPlan.Columns().GatewayProductStatus: res.GatewayPlanStatus,
			dao.GatewayPlan.Columns().Data:                 res.Data,
			dao.GatewayPlan.Columns().Status:               int(res.Status),
		}).Where(dao.GatewayPlan.Columns().Id, gatewayPlan.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	return gerror.Newf("SubscriptionPlanChannelTransferAndActivate update err:%s", update)
		//}
		gatewayPlan.GatewayPlanId = res.GatewayPlanId
		gatewayPlan.GatewayProductStatus = res.GatewayPlanStatus
		gatewayPlan.Data = res.Data
		gatewayPlan.Status = int(res.Status)
	}

	return nil
}
