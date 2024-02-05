package service

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/gateway/api"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strings"
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

func SubscriptionPlanChannelTransferAndActivate(ctx context.Context, planId int64, channelId int64) error {
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
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, channelId)
	utility.Assert(payChannel != nil, "payChannel not found")
	planChannel := query.GetPlanChannel(ctx, planId, channelId)
	if planChannel == nil {
		planChannel = &entity.GatewayPlan{
			PlanId:    planId,
			GatewayId: channelId,
			Status:    consts.PlanChannelStatusInit,
		}
		//保存planChannel
		result, err := dao.GatewayPlan.Ctx(ctx).Data(planChannel).OmitNil().Insert(planChannel)
		if err != nil {
			err = gerror.Newf(`SubscriptionPlanChannelTransferAndActivate record insert failure %s`, err)
			planChannel = nil
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			planChannel = nil
			return err
		}
		planChannel.Id = uint64(uint(id))
	}
	if len(planChannel.GatewayProductId) == 0 {
		//产品尚未创建
		if len(plan.ChannelProductName) == 0 {
			plan.ChannelProductName = plan.PlanName
		}
		if len(plan.ChannelProductDescription) == 0 {
			plan.ChannelProductDescription = plan.Description
		}
		res, err := api.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelProductCreate(ctx, plan, planChannel)
		if err != nil {
			return err
		}
		//更新 planChannel
		_, err = dao.GatewayPlan.Ctx(ctx).Data(g.Map{
			dao.GatewayPlan.Columns().GatewayProductId:     res.GatewayProductId,
			dao.GatewayPlan.Columns().GatewayProductStatus: res.ChannelProductStatus,
		}).Where(dao.GatewayPlan.Columns().Id, planChannel.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	return gerror.Newf("SubscriptionPlanChannelTransferAndActivate update err:%s", update)
		//}
		planChannel.GatewayProductId = res.GatewayProductId
		planChannel.GatewayProductStatus = res.ChannelProductStatus
	}
	if len(planChannel.GatewayPlanId) == 0 {
		//创建 并激活 Plan
		res, err := api.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelPlanCreateAndActivate(ctx, plan, planChannel)
		if err != nil {
			return err
		}
		_, err = dao.GatewayPlan.Ctx(ctx).Data(g.Map{
			dao.GatewayPlan.Columns().GatewayPlanId:        res.GatewayPlanId,
			dao.GatewayPlan.Columns().GatewayProductStatus: res.ChannelPlanStatus,
			dao.GatewayPlan.Columns().Data:                 res.Data,
			dao.GatewayPlan.Columns().Status:               int(res.Status),
		}).Where(dao.GatewayPlan.Columns().Id, planChannel.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	return gerror.Newf("SubscriptionPlanChannelTransferAndActivate update err:%s", update)
		//}
		planChannel.GatewayPlanId = res.GatewayPlanId
		planChannel.GatewayProductStatus = res.ChannelPlanStatus
		planChannel.Data = res.Data
		planChannel.Status = int(res.Status)
	}

	return nil
}
