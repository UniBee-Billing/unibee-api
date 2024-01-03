package service

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "go-oversea-pay/api/subscription/v1"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/payment/outchannel"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func SubscriptionCreate(ctx context.Context, req *v1.SubscriptionCreateReq) (*entity.Subscription, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.PlanId > 0, "PlanId invalid")
	utility.Assert(req.ChannelId > 0, "ConfirmChannelId invalid")
	utility.Assert(req.UserId > 0, "UserId invalid")
	plan := query.GetSubscriptionPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	planChannel := query.GetSubscriptionPlanChannel(ctx, req.PlanId, req.ChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "plan channel should be transfer first")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, req.ChannelId)
	utility.Assert(payChannel != nil, "payChannel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")

	//todo mark plan 是否活跃检查

	one := &entity.Subscription{
		MerchantId: merchantInfo.Id,
		PlanId:     req.PlanId,
		ChannelId:  req.ChannelId,
		UserId:     req.UserId,
		Quantity:   12, // todo mark 按照逻辑计算数量
		//CustomerName:          "jack3",             // todo mark
		//CustomerEmail:         "jack3.fu@wowow.io", // todo mark
		SubscriptionId:        utility.CreateSubscriptionOrderNo(),
		ChannelSubscriptionId: "",
		Status:                consts.SubStatusInit,
		ChannelUserId:         req.ChannelUserId,
		Data:                  "", //额外参数配置
	}

	result, err := dao.Subscription.Ctx(ctx).Data(one).OmitEmpty().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionCreate record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	createRes, err := outchannel.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionCreate(ctx, &ro.ChannelCreateSubscriptionInternalReq{
		Plan:         plan,
		SubPlans:     nil,
		PlanChannel:  planChannel,
		Subscription: one,
	})
	if err != nil {
		return nil, err
	}
	//更新 Subscription
	update, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().ChannelUserId:         createRes.ChannelUserId,
		dao.Subscription.Columns().ChannelSubscriptionId: createRes.ChannelSubscriptionId,
		dao.Subscription.Columns().Status:                consts.SubStatusCreate,
		dao.Subscription.Columns().Link:                  createRes.Link,
		dao.Subscription.Columns().ResponseData:          createRes.Data,
	}).Where(dao.Subscription.Columns().Id, one.Id).OmitEmpty().Update()
	if err != nil {
		return nil, err
	}
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return nil, gerror.Newf("SubscriptionCreate update err:%s", update)
	}
	one.ChannelSubscriptionId = createRes.ChannelSubscriptionId
	one.Status = consts.PlanChannelStatusCreate
	one.Link = createRes.Link
	one.ChannelUserId = createRes.ChannelUserId

	return one, nil
}

func SubscriptionUpdate(ctx context.Context, req *v1.SubscriptionUpdateReq) (*entity.Subscription, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.NewPlanId > 0, "NewPlanId invalid")
	utility.Assert(req.ConfirmChannelId > 0, "ConfirmChannelId invalid")
	utility.Assert(req.SubscriptionId > 0, "SubscriptionId invalid")
	plan := query.GetSubscriptionPlanById(ctx, req.NewPlanId)
	utility.Assert(plan != nil, "invalid planId")
	planChannel := query.GetSubscriptionPlanChannel(ctx, req.NewPlanId, req.ConfirmChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "plan channel should be transfer first")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, req.ConfirmChannelId)
	utility.Assert(payChannel != nil, "payChannel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	subscription := query.GetSubscriptionById(ctx, req.SubscriptionId)
	utility.Assert(subscription != nil, "subscription not found")
	utility.Assert(subscription.ChannelId == req.ConfirmChannelId, "channel not match")
	//暂时不开放不同通道升级功能 todo mark
	oldPlan := query.GetSubscriptionPlanById(ctx, subscription.PlanId)
	utility.Assert(oldPlan != nil, "oldPlan not found")
	oldPlanChannel := query.GetSubscriptionPlanChannel(ctx, int64(oldPlan.Id), req.ConfirmChannelId)
	utility.Assert(oldPlanChannel != nil, "oldPlanChannel not found")

	//todo mark subscription 检查

	//todo mark plan 是否活跃检查
	updateRes, err := outchannel.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionUpdate(ctx, &ro.ChannelUpdateSubscriptionInternalReq{
		Plan:           plan,
		OldPlan:        oldPlan,
		SubPlans:       nil,
		PlanChannel:    planChannel,
		OldPlanChannel: oldPlanChannel,
		Subscription:   subscription,
	})
	if err != nil {
		return nil, err
	}

	//更新 Subscription
	update, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().PlanId:       plan.Id,
		dao.Subscription.Columns().ResponseData: updateRes.Data,
		dao.Subscription.Columns().GmtModify:    gtime.Now(),
		dao.Subscription.Columns().Link:         updateRes.Link,
	}).Where(dao.Subscription.Columns().Id, subscription.Id).OmitEmpty().Update()
	if err != nil {
		return nil, err
	}
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return nil, gerror.Newf("SubscriptionUpdate update err:%s", update)
	}
	subscription.ChannelSubscriptionId = updateRes.ChannelSubscriptionId
	subscription.Status = consts.PlanChannelStatusCreate
	subscription.Link = updateRes.Link

	return subscription, nil
}
