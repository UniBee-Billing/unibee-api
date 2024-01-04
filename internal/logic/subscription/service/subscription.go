package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/api/user/subscription"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/payment/outchannel"
	outchannelro "go-oversea-pay/internal/logic/payment/outchannel/ro"
	"go-oversea-pay/internal/logic/subscription/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strings"
)

type SubscriptionCreatePrepareInternalRes struct {
	Plan         *entity.SubscriptionPlan           `json:"planId"`
	Quantity     int64                              `json:"quantity"`
	PlanChannel  *entity.SubscriptionPlanChannel    `json:"planChannel"`
	PayChannel   *entity.OverseaPayChannel          `json:"payChannel"`
	MerchantInfo *entity.MerchantInfo               `json:"merchantInfo"`
	AddonParams  []*ro.SubscriptionPlanAddonParamRo `json:"addonParams"`
	Addons       []*ro.SubscriptionPlanAddonRo      `json:"addons"`
	TotalAmount  int64                              `json:"totalAmount"                ` // 金额,单位：分
	Currency     string                             `json:"currency"              `      // 货币
	UserId       int64                              `json:"UserId" `
}

func SubscriptionCreatePrepare(ctx context.Context, req *subscription.SubscriptionCreatePrepareReq) (*SubscriptionCreatePrepareInternalRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.PlanId > 0, "PlanId invalid")
	utility.Assert(req.ChannelId > 0, "ConfirmChannelId invalid")
	utility.Assert(req.UserId > 0, "UserId invalid")
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(int64(_interface.BizCtx().Get(ctx).User.Id) == req.UserId, "userId not match")
	}
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	planChannel := query.GetPlanChannel(ctx, req.PlanId, req.ChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "plan channel should be transfer first")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, req.ChannelId)
	utility.Assert(payChannel != nil, "payChannel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")

	//设置默认值
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	var totalAmount = plan.Amount * int64(req.Quantity)
	var currency = plan.Currency

	addons := checkAndListAddonsFromParams(ctx, req.AddonParams)

	for _, addon := range addons {
		utility.Assert(strings.Compare(addon.AddonPlan.Currency, currency) == 0, fmt.Sprintf("currency not match for planId:%v addonId:%v", plan.Id, addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.MerchantId == plan.MerchantId, fmt.Sprintf("Addon Id:%v Merchant not match", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.Status == consts.PlanStatusPublished, fmt.Sprintf("Addon Id:%v Not Publish status", addon.AddonPlan.Id))
		totalAmount = totalAmount + addon.AddonPlan.Amount*int64(addon.Quantity)
	}

	return &SubscriptionCreatePrepareInternalRes{
		Plan:         plan,
		PlanChannel:  planChannel,
		PayChannel:   payChannel,
		MerchantInfo: merchantInfo,
		AddonParams:  req.AddonParams,
		Addons:       addons,
		TotalAmount:  totalAmount,
		Currency:     currency,
		UserId:       req.UserId,
	}, nil
}

func checkAndListAddonsFromParams(ctx context.Context, addonParams []*ro.SubscriptionPlanAddonParamRo) []*ro.SubscriptionPlanAddonRo {
	var addons []*ro.SubscriptionPlanAddonRo
	var totalAddonIds []int64
	if len(addonParams) > 0 {
		for _, s := range addonParams {
			totalAddonIds = append(totalAddonIds, s.AddonPlanId) // 添加到整数列表中
		}
	}
	var allAddonList []*entity.SubscriptionPlan
	if len(totalAddonIds) > 0 {
		//查询所有 Plan
		err := dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, totalAddonIds).Scan(&allAddonList)
		if err == nil {
			//整合进列表
			mapPlans := make(map[int64]*entity.SubscriptionPlan)
			for _, pair := range allAddonList {
				key := int64(pair.Id)
				value := pair
				mapPlans[key] = value
			}
			for _, param := range addonParams {
				//所有 Addon 项目必须要能查到
				//类型是 Addon
				//未删除
				//数量大于 0
				utility.Assert(mapPlans[param.AddonPlanId] != nil, fmt.Sprintf("AddonPlanId not found:%v", param.AddonPlanId))
				utility.Assert(mapPlans[param.AddonPlanId].Type == consts.PlanTypeAddon, fmt.Sprintf("Id:%v not Addon Type", param.AddonPlanId))
				utility.Assert(mapPlans[param.AddonPlanId].IsDeleted == 0, fmt.Sprintf("Addon Id:%v is Deleted", param.AddonPlanId))
				utility.Assert(param.Quantity > 0, fmt.Sprintf("Id:%v quantity invalid", param.AddonPlanId))
				addons = append(addons, &ro.SubscriptionPlanAddonRo{
					Quantity:  param.Quantity,
					AddonPlan: mapPlans[param.AddonPlanId],
				})
			}
		}
	}
	return addons
}

func SubscriptionCreate(ctx context.Context, req *subscription.SubscriptionCreateReq) (*entity.Subscription, error) {
	prepare, err := SubscriptionCreatePrepare(ctx, &subscription.SubscriptionCreatePrepareReq{
		PlanId:      req.PlanId,
		Quantity:    req.Quantity,
		ChannelId:   req.ChannelId,
		UserId:      req.UserId,
		AddonParams: req.AddonParams,
	})
	if err != nil {
		return nil, err
	}

	//校验
	utility.Assert(req.ConfirmTotalAmount == prepare.TotalAmount, "totalAmount not match , data may expired, fetch again")
	utility.Assert(strings.Compare(req.ConfirmCurrency, prepare.Currency) == 0, "currency not match , data may expired, fetch again")
	//channelUserId 处理
	var channelUserId string
	channelUser := query.GetUserChannel(ctx, prepare.UserId, prepare.PlanChannel.ChannelId)
	if channelUser != nil {
		channelUserId = channelUser.ChannelUserId
	}

	one := &entity.Subscription{
		MerchantId:            prepare.MerchantInfo.Id,
		PlanId:                int64(prepare.Plan.Id),
		ChannelId:             prepare.PlanChannel.ChannelId,
		UserId:                prepare.UserId,
		Quantity:              prepare.Quantity,
		Amount:                prepare.TotalAmount,
		Currency:              prepare.Currency,
		AddonData:             utility.MarshalToJsonString(prepare.Addons),
		SubscriptionId:        utility.CreateSubscriptionOrderNo(),
		ChannelSubscriptionId: "",
		Status:                consts.SubStatusInit,
		ChannelUserId:         channelUserId,
		Data:                  "", //额外参数配置
	}

	result, err := dao.Subscription.Ctx(ctx).Data(one).OmitEmpty().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionCreate record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	createRes, err := outchannel.GetPayChannelServiceProvider(ctx, int64(prepare.PayChannel.Id)).DoRemoteChannelSubscriptionCreate(ctx, &outchannelro.ChannelCreateSubscriptionInternalReq{
		Plan:         prepare.Plan,
		AddonPlans:   prepare.Addons,
		PlanChannel:  prepare.PlanChannel,
		Subscription: one,
	})
	if err != nil {
		return nil, err
	}

	//更新 Subscription
	update, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().ChannelUserId:         createRes.ChannelUserId, // todo mark 进行 UserId 匹配更新
		dao.Subscription.Columns().ChannelSubscriptionId: createRes.ChannelSubscriptionId,
		dao.Subscription.Columns().Status:                consts.SubStatusCreate,
		dao.Subscription.Columns().Link:                  createRes.Link,
		dao.Subscription.Columns().ResponseData:          createRes.Data,
		dao.Subscription.Columns().GmtModify:             gtime.Now(),
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

	if channelUser == nil && len(createRes.ChannelUserId) > 0 {
		_, err := query.SaveUserChannel(ctx, prepare.UserId, prepare.PlanChannel.ChannelId, channelUserId)
		if err != nil {
			// ChannelUser 创建错误
			return nil, gerror.Newf("SubscriptionCreate ChannelUser save err:%s", err)
		}
	}

	return one, nil
}

func SubscriptionUpdate(ctx context.Context, req *subscription.SubscriptionUpdateReq) (*entity.SubscriptionPendingUpdate, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.NewPlanId > 0, "NewPlanId invalid")
	utility.Assert(req.ConfirmChannelId > 0, "ConfirmChannelId invalid")
	utility.Assert(len(req.SubscriptionId) > 0, "SubscriptionId invalid")
	plan := query.GetPlanById(ctx, req.NewPlanId)
	utility.Assert(plan != nil, "invalid planId")
	planChannel := query.GetPlanChannel(ctx, req.NewPlanId, req.ConfirmChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "plan channel should be transfer first")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, req.ConfirmChannelId)
	utility.Assert(payChannel != nil, "payChannel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	subscription := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(subscription != nil, "subscription not found")
	utility.Assert(subscription.ChannelId == req.ConfirmChannelId, "channel not match")
	//暂时不开放不同通道升级功能 todo mark
	oldPlan := query.GetPlanById(ctx, subscription.PlanId)
	utility.Assert(oldPlan != nil, "oldPlan not found")
	oldPlanChannel := query.GetPlanChannel(ctx, int64(oldPlan.Id), req.ConfirmChannelId)
	utility.Assert(oldPlanChannel != nil, "oldPlanChannel not found")

	//todo mark subscription 检查

	one := &entity.SubscriptionPendingUpdate{
		MerchantId:           merchantInfo.Id,
		ChannelId:            subscription.ChannelId,
		UserId:               subscription.UserId,
		SubscriptionId:       subscription.SubscriptionId,
		UpdateSubscriptionId: utility.CreateSubscriptionOrderNo(),
		Amount:               subscription.Amount,
		Currency:             subscription.Currency,
		PlanId:               subscription.PlanId,
		Quantity:             subscription.Quantity,
		AddonData:            subscription.AddonData,
		UpdateAmount:         subscription.Amount, //总金额 todo mark 需要添加 Addon，并用计算函数重新计算
		UpdateCurrency:       plan.Currency,
		UpdatePlanId:         req.NewPlanId,
		UpdateQuantity:       1,                      //todo mark 主 plan 暂时不支持数量调整
		UpdatedAddonData:     subscription.AddonData, // addon 带上之前订阅
		Status:               consts.SubStatusInit,
		Data:                 "", //额外参数配置
	}

	result, err := dao.SubscriptionPendingUpdate.Ctx(ctx).Data(one).OmitEmpty().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionPendingUpdate record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	//todo mark plan 是否活跃检查
	updateRes, err := outchannel.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionUpdate(ctx, &outchannelro.ChannelUpdateSubscriptionInternalReq{
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
	update, err := dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPendingUpdate.Columns().ResponseData: updateRes.Data,
		dao.SubscriptionPendingUpdate.Columns().GmtModify:    gtime.Now(),
		dao.SubscriptionPendingUpdate.Columns().Link:         updateRes.Link,
	}).Where(dao.SubscriptionPendingUpdate.Columns().Id, one.Id).OmitEmpty().Update()
	if err != nil {
		return nil, err
	}
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return nil, gerror.Newf("SubscriptionUpdate update err:%s", update)
	}
	one.ChannelUpdateId = updateRes.ChannelSubscriptionId
	one.Status = consts.PlanChannelStatusCreate
	one.Link = updateRes.Link

	return one, nil
}
