package service

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/channel/out"
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
		planChannel = &entity.ChannelPlan{
			PlanId:    planId,
			ChannelId: channelId,
			Status:    consts.PlanChannelStatusInit,
		}
		//保存planChannel
		result, err := dao.ChannelPlan.Ctx(ctx).Data(planChannel).OmitNil().Insert(planChannel)
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
	if len(planChannel.ChannelProductId) == 0 {
		//产品尚未创建
		if len(plan.ChannelProductName) == 0 {
			plan.ChannelProductName = plan.PlanName
		}
		if len(plan.ChannelProductDescription) == 0 {
			plan.ChannelProductDescription = plan.Description
		}
		res, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelProductCreate(ctx, plan, planChannel)
		if err != nil {
			return err
		}
		//更新 planChannel
		_, err = dao.ChannelPlan.Ctx(ctx).Data(g.Map{
			dao.ChannelPlan.Columns().ChannelProductId:     res.ChannelProductId,
			dao.ChannelPlan.Columns().ChannelProductStatus: res.ChannelProductStatus,
		}).Where(dao.ChannelPlan.Columns().Id, planChannel.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	return gerror.Newf("SubscriptionPlanChannelTransferAndActivate update err:%s", update)
		//}
		planChannel.ChannelProductId = res.ChannelProductId
		planChannel.ChannelProductStatus = res.ChannelProductStatus
	}
	if len(planChannel.ChannelPlanId) == 0 {
		//创建 并激活 Plan
		res, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelPlanCreateAndActivate(ctx, plan, planChannel)
		if err != nil {
			return err
		}
		_, err = dao.ChannelPlan.Ctx(ctx).Data(g.Map{
			dao.ChannelPlan.Columns().ChannelPlanId:     res.ChannelPlanId,
			dao.ChannelPlan.Columns().ChannelPlanStatus: res.ChannelPlanStatus,
			dao.ChannelPlan.Columns().Data:              res.Data,
			dao.ChannelPlan.Columns().Status:            int(res.Status),
		}).Where(dao.ChannelPlan.Columns().Id, planChannel.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	return gerror.Newf("SubscriptionPlanChannelTransferAndActivate update err:%s", update)
		//}
		planChannel.ChannelPlanId = res.ChannelPlanId
		planChannel.ChannelPlanStatus = res.ChannelPlanStatus
		planChannel.Data = res.Data
		planChannel.Status = int(res.Status)
	}

	return nil
}
