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
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strings"
)

func SubscriptionPlanChannelActivate(ctx context.Context, planId int64, channelId int64) (err error) {
	utility.Assert(planId > 0, "invalid planId")
	utility.Assert(channelId > 0, "invalid channelId")
	plan := query.GetSubscriptionPlanById(ctx, planId)
	utility.Assert(plan != nil, "invalid planId")
	planChannel := query.GetSubscriptionPlanChannel(ctx, planId, channelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "plan channel should be transfer first")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, channelId)
	utility.Assert(payChannel != nil, "payChannel not found")

	err = outchannel.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelPlanActive(ctx, plan, planChannel)
	if err != nil {
		return
	}
	update, err := dao.SubscriptionPlanChannel.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPlanChannel.Columns().Status: consts.PlanStatusActive,
		//dao.SubscriptionPlanChannel.Columns().ChannelPlanStatus: consts.PlanStatusActive,// todo mark
		dao.SubscriptionPlanChannel.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionPlanChannel.Columns().Id, planChannel.Id).Update()
	if err != nil {
		return err
	}
	// todo mark update 值没变化会报错
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return gerror.Newf("SubscriptionPlanChannelActivate update err:%s", update)
	}
	return
}

func SubscriptionPlanChannelDeactivate(ctx context.Context, planId int64, channelId int64) (err error) {
	utility.Assert(planId > 0, "invalid planId")
	utility.Assert(channelId > 0, "invalid channelId")
	plan := query.GetSubscriptionPlanById(ctx, planId)
	utility.Assert(plan != nil, "invalid planId")
	planChannel := query.GetSubscriptionPlanChannel(ctx, planId, channelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "plan channel should be transfer first")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, channelId)
	utility.Assert(payChannel != nil, "payChannel not found")

	err = outchannel.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelPlanDeactivate(ctx, plan, planChannel)
	if err != nil {
		return
	}
	update, err := dao.SubscriptionPlanChannel.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPlanChannel.Columns().Status: consts.PlanStatusInActive,
		//dao.SubscriptionPlanChannel.Columns().ChannelPlanStatus: consts.PlanStatusInActive,// todo mark
		dao.SubscriptionPlanChannel.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionPlanChannel.Columns().Id, planChannel.Id).Update()
	if err != nil {
		return err
	}
	// todo mark update 值没变化会报错
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return gerror.Newf("SubscriptionPlanChannelDeactivate update err:%s", update)
	}
	return
}

func SubscriptionPlanCreate(ctx context.Context, req *v1.SubscriptionPlanCreateReq) (one *entity.SubscriptionPlan, err error) {
	intervals := []string{"day", "month", "year", "week"}
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.Amount > 0, "amount value should > 0")
	utility.Assert(len(req.ImageUrl) > 0, "imageUrl should not be null")
	utility.Assert(strings.HasPrefix(req.ImageUrl, "http"), "imageUrl should start with http")
	merchantInfo := query.GetMerchantInfoById(ctx, req.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	utility.Assert(utility.StringContainsElement(intervals, strings.ToLower(req.IntervalUnit)), "IntervalUnit 错误，day｜month｜year｜week\"")
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
		Amount:                    req.Amount,
		Currency:                  strings.ToUpper(req.Currency),
		IntervalUnit:              strings.ToLower(req.IntervalUnit),
		Description:               req.Description,
		ImageUrl:                  req.ImageUrl,
		HomeUrl:                   req.HomeUrl,
		ChannelProductName:        req.ProductName,
		ChannelProductDescription: req.ProductDescription,
	}
	result, err := dao.SubscriptionPlan.Ctx(ctx).Data(one).OmitEmpty().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionPlanCreate record insert failure %s`, err)
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
		res, err := outchannel.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelProductCreate(ctx, plan, planChannel)
		if err != nil {
			return err
		}
		//更新 planChannel
		update, err := dao.SubscriptionPlanChannel.Ctx(ctx).Data(g.Map{
			dao.SubscriptionPlanChannel.Columns().ChannelProductId:     res.ChannelProductId,
			dao.SubscriptionPlanChannel.Columns().ChannelProductStatus: res.ChannelProductStatus,
		}).Where(dao.SubscriptionPlanChannel.Columns().Id, planChannel.Id).Update()
		if err != nil {
			return err
		}
		rowAffected, err := update.RowsAffected()
		if rowAffected != 1 {
			return gerror.Newf("SubscriptionPlanChannelTransferAndActivate update err:%s", update)
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
		update, err := dao.SubscriptionPlanChannel.Ctx(ctx).Data(g.Map{
			dao.SubscriptionPlanChannel.Columns().ChannelPlanId:     res.ChannelPlanId,
			dao.SubscriptionPlanChannel.Columns().ChannelPlanStatus: res.ChannelPlanStatus,
			dao.SubscriptionPlanChannel.Columns().Data:              res.Data,
			dao.SubscriptionPlanChannel.Columns().Status:            int(res.Status),
		}).Where(dao.SubscriptionPlanChannel.Columns().Id, planChannel.Id).Update()
		if err != nil {
			return err
		}
		rowAffected, err := update.RowsAffected()
		if rowAffected != 1 {
			return gerror.Newf("SubscriptionPlanChannelTransferAndActivate update err:%s", update)
		}
		planChannel.ChannelPlanId = res.ChannelPlanId
		planChannel.ChannelPlanStatus = res.ChannelPlanStatus
		planChannel.Data = res.Data
		planChannel.Status = int(res.Status)
	}

	return nil
}
