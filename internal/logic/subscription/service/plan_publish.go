package service

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/gateway"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
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
	plan := query.GetPlanById(ctx, planId)
	utility.Assert(plan != nil, "plan not found")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, channelId)
	utility.Assert(payChannel != nil, "payChannel not found")
	planChannel := query.GetPlanChannel(ctx, planId, channelId)
	if planChannel == nil {
		planChannel = &entity.SubscriptionPlanChannel{
			PlanId:    planId,
			ChannelId: channelId,
			Status:    consts.PlanChannelStatusInit,
		}
		//保存planChannel
		result, err := dao.SubscriptionPlanChannel.Ctx(ctx).Data(planChannel).OmitNil().Insert(planChannel)
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
		res, err := gateway.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelProductCreate(ctx, plan, planChannel)
		if err != nil {
			return err
		}
		//更新 planChannel
		_, err = dao.SubscriptionPlanChannel.Ctx(ctx).Data(g.Map{
			dao.SubscriptionPlanChannel.Columns().ChannelProductId:     res.ChannelProductId,
			dao.SubscriptionPlanChannel.Columns().ChannelProductStatus: res.ChannelProductStatus,
		}).Where(dao.SubscriptionPlanChannel.Columns().Id, planChannel.Id).OmitNil().Update()
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
		res, err := gateway.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelPlanCreateAndActivate(ctx, plan, planChannel)
		if err != nil {
			return err
		}
		_, err = dao.SubscriptionPlanChannel.Ctx(ctx).Data(g.Map{
			dao.SubscriptionPlanChannel.Columns().ChannelPlanId:     res.ChannelPlanId,
			dao.SubscriptionPlanChannel.Columns().ChannelPlanStatus: res.ChannelPlanStatus,
			dao.SubscriptionPlanChannel.Columns().Data:              res.Data,
			dao.SubscriptionPlanChannel.Columns().Status:            int(res.Status),
		}).Where(dao.SubscriptionPlanChannel.Columns().Id, planChannel.Id).OmitNil().Update()
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
