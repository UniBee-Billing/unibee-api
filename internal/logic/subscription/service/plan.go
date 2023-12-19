package service

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	v1 "go-oversea-pay/api/subscription/v1"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/payment/outchannel"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strconv"
)

func SubscriptionPlanCreate(ctx context.Context, req *v1.SubscriptionPlanCreateReq) (one *entity.SubscriptionPlan, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.Amount > 0, "amount value should > 0")
	merchantInfo := query.GetMerchantInfoById(ctx, req.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	if len(req.ProductName) == 0 {
		req.ProductName = req.PlanName
	}
	if len(req.ProductDescription) == 0 {
		req.ProductDescription = req.Description
	}
	one = &entity.SubscriptionPlan{
		CompanyId:                 merchantInfo.CompanyId,
		MerchantId:                req.MerchantId,
		PlanName:                  req.PlanName,
		Amount:                    strconv.FormatInt(req.Amount, 10),
		Currency:                  req.Currency,
		IntervalUnit:              req.IntervalUnit,
		Description:               req.Description,
		ImageUrl:                  req.ImageUrl,
		HomeUrl:                   req.HomeUrl,
		ChannelProductName:        req.ProductName,
		ChannelProductDescription: req.ProductDescription,
	}
	result, err := dao.SubscriptionPlan.Ctx(ctx).Data(one).OmitEmpty().Insert(one)
	if err != nil {
		err = gerror.Newf(`record insert failure %s`, err)
		one = nil
		return
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))
	return one, nil
}

func SubscriptionPlanChannelTransferAndActivate(ctx context.Context, planId int64, channelId int64) error {
	plan := query.GetSubscriptionPlanById(ctx, planId)
	utility.Assert(plan != nil, "plan not found")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, channelId)
	utility.Assert(payChannel != nil, "payChannel not found")
	planChannel := query.GetSubscriptionPlanChannel(ctx, planId, channelId)
	if planChannel == nil {
		planChannel = &entity.SubscriptionPlanChannel{
			PlanId:    planId,
			ChannelId: channelId,
			Status:    consts.PlanStatusInit,
		}
		//保存planChannel
		result, err := dao.SubscriptionPlanChannel.Ctx(ctx).Data(planChannel).OmitEmpty().Insert(planChannel)
		if err != nil {
			err = gerror.Newf(`record insert failure %s`, err)
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
		res, err := outchannel.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelProductCreate(ctx, plan, planChannel)
		if err != nil {
			return err
		}
		//更新 planChannel
		update, err := dao.SubscriptionPlanChannel.Ctx(ctx).Where(entity.SubscriptionPlanChannel{Id: planChannel.Id}).Update(entity.SubscriptionPlanChannel{
			ChannelProductId:     res.ChannelProductId,
			ChannelProductStatus: res.ChannelProductStatus,
		})
		if err != nil {
			return err
		}
		rowAffected, err := update.RowsAffected()
		if rowAffected != 1 {
			return gerror.Newf("update err:%s", update)
		}
		planChannel.ChannelProductId = res.ChannelProductId
		planChannel.ChannelProductStatus = res.ChannelProductStatus
	}
	if len(planChannel.ChannelPlanId) == 0 {
		//创建 并激活 Plan
		res, err := outchannel.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelPlanCreateAndActivate(ctx, plan, planChannel)
		if err != nil {
			return err
		}
		update, err := dao.SubscriptionPlanChannel.Ctx(ctx).Where(entity.SubscriptionPlanChannel{Id: planChannel.Id}).Update(entity.SubscriptionPlanChannel{
			ChannelPlanId:     res.ChannelPlanId,
			ChannelPlanStatus: res.ChannelPlanStatus,
			Data:              res.Data,
			Status:            int(res.Status),
		})
		if err != nil {
			return err
		}
		rowAffected, err := update.RowsAffected()
		if rowAffected != 1 {
			return gerror.Newf("update err:%s", update)
		}
		planChannel.ChannelPlanId = res.ChannelPlanId
		planChannel.ChannelPlanStatus = res.ChannelPlanStatus
		planChannel.Data = res.Data
		planChannel.Status = int(res.Status)
	}

	return nil
}
