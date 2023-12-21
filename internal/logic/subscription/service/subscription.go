package service

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	v1 "go-oversea-pay/api/subscription/v1"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/payment/outchannel"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func SubscriptionCreate(ctx context.Context, req *v1.SubscriptionCreateReq) (res *entity.Subscription, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.PlanId > 0, "PlanId invalid")
	utility.Assert(req.ChannelId > 0, "ChannelId invalid")
	utility.Assert(req.UserId > 0, "UserId invalid")
	plan := query.GetSubscriptionPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	planChannel := query.GetSubscriptionPlanChannel(ctx, req.PlanId, req.ChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "plan channel should be transfer first")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, req.ChannelId)
	utility.Assert(payChannel != nil, "payChannel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")

	one := &entity.Subscription{
		CompanyId:             merchantInfo.CompanyId,
		MerchantId:            merchantInfo.Id,
		PlanId:                req.PlanId,
		ChannelId:             req.ChannelId,
		UserId:                req.UserId,
		Quantity:              12, // todo mark 按照逻辑计算数量
		SubscriptionId:        utility.CreateSubscriptionOrderNo(),
		ChannelSubscriptionId: "",
		Status:                consts.SubStatusInit,
		ChannelUserId:         req.ChannelUserId,
	}

	result, err := dao.Subscription.Ctx(ctx).Data(one).OmitEmpty().Insert(one)
	if err != nil {
		err = gerror.Newf(`record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	createRes, err := outchannel.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionCreate(ctx, plan, planChannel, one)
	if err != nil {
		return nil, err
	}
	//更新 planChannel
	update, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().ChannelSubscriptionId: createRes.ChannelSubscriptionId,
		dao.Subscription.Columns().Status:                consts.SubStatusCreate, //todo mark createRes 判断状态
	}).Where(dao.Subscription.Columns().Id, one.Id).Update()
	if err != nil {
		return nil, err
	}
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return nil, gerror.Newf("update err:%s", update)
	}
	one.ChannelSubscriptionId = createRes.ChannelSubscriptionId

	return one, nil
}
