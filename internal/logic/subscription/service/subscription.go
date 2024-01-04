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
	Invoice      *ro.SubscriptionInvoiceRo          `json:"invoice"`
	UserId       int64                              `json:"userId" `
	Email        string                             `json:"email" `
}

func checkAndListAddonsFromParams(ctx context.Context, addonParams []*ro.SubscriptionPlanAddonParamRo, channelId int64) []*ro.SubscriptionPlanAddonRo {
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
				planChannel := query.GetPlanChannel(ctx, int64(mapPlans[param.AddonPlanId].Id), channelId) // todo mark for 循环内调用 需做缓存，此数据基本不会变化,或者方案 2 使用 channelId 合并查询
				utility.Assert(len(planChannel.ChannelPlanId) > 0, fmt.Sprintf("internal error PlanId:%v ChannelId:%v channelPlanId invalid", param.AddonPlanId, channelId))
				utility.Assert(planChannel.Status == consts.PlanChannelStatusActive, fmt.Sprintf("internal error PlanId:%v ChannelId:%v channelPlanStatus not active", param.AddonPlanId, channelId))
				addons = append(addons, &ro.SubscriptionPlanAddonRo{
					Quantity:         param.Quantity,
					AddonPlan:        mapPlans[param.AddonPlanId],
					AddonPlanChannel: planChannel,
				})
			}
		}
	}
	return addons
}

func SubscriptionCreatePreview(ctx context.Context, req *subscription.SubscriptionCreatePreviewReq) (*SubscriptionCreatePrepareInternalRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.PlanId > 0, "PlanId invalid")
	utility.Assert(req.ChannelId > 0, "ConfirmChannelId invalid")
	utility.Assert(req.UserId > 0, "UserId invalid")
	email := ""
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(int64(_interface.BizCtx().Get(ctx).User.Id) == req.UserId, "userId not match")
		email = _interface.BizCtx().Get(ctx).User.Email
	}
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusPublished, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	planChannel := query.GetPlanChannel(ctx, req.PlanId, req.ChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "internal error plan channel transfer not complete")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, req.ChannelId) //todo mark 改造成支持 Merchant 级别的 PayChannel
	utility.Assert(payChannel != nil, "payChannel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")

	//设置默认值
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	var totalAmount = plan.Amount * req.Quantity
	var currency = plan.Currency

	addons := checkAndListAddonsFromParams(ctx, req.AddonParams, planChannel.ChannelId)

	for _, addon := range addons {
		utility.Assert(strings.Compare(addon.AddonPlan.Currency, currency) == 0, fmt.Sprintf("currency not match for planId:%v addonId:%v", plan.Id, addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.MerchantId == plan.MerchantId, fmt.Sprintf("Addon Id:%v Merchant not match", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.Status == consts.PlanStatusPublished, fmt.Sprintf("Addon Id:%v Not Publish status", addon.AddonPlan.Id))
		totalAmount = totalAmount + addon.AddonPlan.Amount*addon.Quantity
	}

	//生成临时账单
	var invoiceItems []*ro.SubscriptionInvoiceItemRo
	invoiceItems = append(invoiceItems, &ro.SubscriptionInvoiceItemRo{
		Currency:    currency,
		Amount:      req.Quantity * plan.Amount,
		Description: plan.PlanName,
	})
	for _, addon := range addons {
		invoiceItems = append(invoiceItems, &ro.SubscriptionInvoiceItemRo{
			Currency:    currency,
			Amount:      addon.Quantity * addon.AddonPlan.Amount,
			Description: addon.AddonPlan.PlanName,
		})
	}

	return &SubscriptionCreatePrepareInternalRes{
		Plan:         plan,
		Quantity:     req.Quantity,
		PlanChannel:  planChannel,
		PayChannel:   payChannel,
		MerchantInfo: merchantInfo,
		AddonParams:  req.AddonParams,
		Addons:       addons,
		TotalAmount:  totalAmount,
		Currency:     currency,
		UserId:       req.UserId,
		Email:        email,
		Invoice: &ro.SubscriptionInvoiceRo{
			TotalAmount:        totalAmount,
			Currency:           currency,
			TaxAmount:          0, // todo mark 暂时不处理 TaxAmount
			SubscriptionAmount: totalAmount,
			Lines:              invoiceItems,
		},
	}, nil
}

func SubscriptionCreate(ctx context.Context, req *subscription.SubscriptionCreateReq) (*entity.Subscription, error) {
	prepare, err := SubscriptionCreatePreview(ctx, &subscription.SubscriptionCreatePreviewReq{
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
		MerchantId:     prepare.MerchantInfo.Id,
		PlanId:         int64(prepare.Plan.Id),
		ChannelId:      prepare.PlanChannel.ChannelId,
		UserId:         prepare.UserId,
		Quantity:       prepare.Quantity,
		Amount:         prepare.TotalAmount,
		Currency:       prepare.Currency,
		AddonData:      utility.MarshalToJsonString(prepare.Addons),
		SubscriptionId: utility.CreateSubscriptionOrderNo(),
		Status:         consts.SubStatusInit,
		CustomerEmail:  prepare.Email,
		ChannelUserId:  channelUserId,
		Data:           "", //额外参数配置
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
		dao.Subscription.Columns().ChannelUserId:         createRes.ChannelUserId,
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

type SubscriptionUpdatePrepareInternalRes struct {
	Subscription   *entity.Subscription               `json:"subscription"`
	Plan           *entity.SubscriptionPlan           `json:"planId"`
	Quantity       int64                              `json:"quantity"`
	PlanChannel    *entity.SubscriptionPlanChannel    `json:"planChannel"`
	PayChannel     *entity.OverseaPayChannel          `json:"payChannel"`
	MerchantInfo   *entity.MerchantInfo               `json:"merchantInfo"`
	AddonParams    []*ro.SubscriptionPlanAddonParamRo `json:"addonParams"`
	Addons         []*ro.SubscriptionPlanAddonRo      `json:"addons"`
	TotalAmount    int64                              `json:"totalAmount"                ` // 金额,单位：分
	Currency       string                             `json:"currency"              `      // 货币
	UserId         int64                              `json:"userId" `
	Email          string                             `json:"email" `
	OldPlan        *entity.SubscriptionPlan           `json:"oldPlan"`
	OldPlanChannel *entity.SubscriptionPlanChannel    `json:"oldPlanChannel"`
	Invoice        *ro.SubscriptionInvoiceRo          `json:"invoice"`
	ProrationDate  int64                              `json:"prorationDate"`
}

func SubscriptionUpdatePreview(ctx context.Context, req *subscription.SubscriptionUpdatePreviewReq) (res *SubscriptionUpdatePrepareInternalRes, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.NewPlanId > 0, "PlanId invalid")
	utility.Assert(len(req.SubscriptionId) > 0, "SubscriptionId invalid")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	//utility.Assert(sub.ChannelId == req.ConfirmChannelId, "channel not match")

	email := ""
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(int64(_interface.BizCtx().Get(ctx).User.Id) == sub.UserId, "userId not match")
		email = _interface.BizCtx().Get(ctx).User.Email
	}
	plan := query.GetPlanById(ctx, req.NewPlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusPublished, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	planChannel := query.GetPlanChannel(ctx, req.NewPlanId, sub.ChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "internal error plan channel transfer not complete")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, sub.ChannelId) //todo mark 改造成支持 Merchant 级别的 PayChannel
	utility.Assert(payChannel != nil, "payChannel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")

	//设置默认值
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	addons := checkAndListAddonsFromParams(ctx, req.AddonParams, planChannel.ChannelId)

	var currency = sub.Currency
	for _, addon := range addons {
		utility.Assert(strings.Compare(addon.AddonPlan.Currency, currency) == 0, fmt.Sprintf("currency not match for planId:%v addonId:%v", plan.Id, addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.MerchantId == plan.MerchantId, fmt.Sprintf("Addon Id:%v Merchant not match", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.Status == consts.PlanStatusPublished, fmt.Sprintf("Addon Id:%v Not Publish status", addon.AddonPlan.Id))
	}

	oldPlan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(oldPlan != nil, "oldPlan not found")
	//暂时不开放不同通道升级功能 todo mark
	oldPlanChannel := query.GetPlanChannel(ctx, int64(oldPlan.Id), sub.ChannelId)
	utility.Assert(oldPlanChannel != nil, "oldPlanChannel not found")
	updatePreviewRes, err := outchannel.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionUpdatePreview(ctx, &outchannelro.ChannelUpdateSubscriptionInternalReq{
		Plan:           plan,
		OldPlan:        oldPlan,
		AddonPlans:     addons, //todo mark oldAddonPlans 是否需要传
		PlanChannel:    planChannel,
		OldPlanChannel: oldPlanChannel,
		Subscription:   sub,
	})
	if err != nil {
		return nil, err
	}
	utility.Assert(strings.Compare(updatePreviewRes.Currency, currency) == 0, fmt.Sprintf("preview currency not match for subscriptionId:%v preview currency:%s", sub.SubscriptionId, updatePreviewRes.Currency))

	var totalAmount = updatePreviewRes.TotalAmount

	return &SubscriptionUpdatePrepareInternalRes{
		Subscription:   sub,
		Plan:           plan,
		Quantity:       req.Quantity,
		PlanChannel:    planChannel,
		PayChannel:     payChannel,
		MerchantInfo:   merchantInfo,
		AddonParams:    req.AddonParams,
		Addons:         addons,
		TotalAmount:    totalAmount,
		Currency:       currency,
		UserId:         sub.UserId,
		Email:          email,
		OldPlan:        oldPlan,
		OldPlanChannel: oldPlanChannel,
		Invoice:        updatePreviewRes.Invoice,
		ProrationDate:  updatePreviewRes.ProrationDate,
	}, nil

}

func SubscriptionUpdate(ctx context.Context, req *subscription.SubscriptionUpdateReq) (*entity.SubscriptionPendingUpdate, error) {
	prepare, err := SubscriptionUpdatePreview(ctx, &subscription.SubscriptionUpdatePreviewReq{
		SubscriptionId: req.SubscriptionId,
		NewPlanId:      req.NewPlanId,
		Quantity:       req.Quantity,
		AddonParams:    req.AddonParams,
	})
	if err != nil {
		return nil, err
	}

	//subscription prepare 检查
	utility.Assert(req.ConfirmTotalAmount == prepare.TotalAmount, "totalAmount not match , data may expired, fetch again")
	utility.Assert(strings.Compare(req.ConfirmCurrency, prepare.Currency) == 0, "currency not match , data may expired, fetch again")

	one := &entity.SubscriptionPendingUpdate{
		MerchantId:           prepare.MerchantInfo.Id,
		ChannelId:            prepare.Subscription.ChannelId,
		UserId:               prepare.Subscription.UserId,
		SubscriptionId:       prepare.Subscription.SubscriptionId,
		UpdateSubscriptionId: utility.CreateSubscriptionOrderNo(),
		Amount:               prepare.Subscription.Amount,
		Currency:             prepare.Subscription.Currency,
		PlanId:               prepare.Subscription.PlanId,
		Quantity:             prepare.Subscription.Quantity,
		AddonData:            prepare.Subscription.AddonData,
		UpdateAmount:         prepare.TotalAmount,
		UpdateCurrency:       prepare.Currency,
		UpdatePlanId:         int64(prepare.Plan.Id),
		UpdateQuantity:       prepare.Quantity,
		UpdatedAddonData:     utility.MarshalToJsonString(prepare.Addons), // addon 暂定不带上之前订阅
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

	updateRes, err := outchannel.GetPayChannelServiceProvider(ctx, int64(prepare.PayChannel.Id)).DoRemoteChannelSubscriptionUpdate(ctx, &outchannelro.ChannelUpdateSubscriptionInternalReq{
		Plan:           prepare.Plan,
		OldPlan:        prepare.OldPlan,
		AddonPlans:     prepare.Addons, //todo mark oldAddonPlans 是否需要传
		PlanChannel:    prepare.PlanChannel,
		OldPlanChannel: prepare.OldPlanChannel,
		Subscription:   prepare.Subscription,
		ProrationDate:  req.ProrationDate,
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
		return nil, gerror.Newf("SubscriptionPendingUpdate update err:%s", update)
	}
	one.ChannelUpdateId = updateRes.ChannelSubscriptionId
	one.Status = consts.PlanChannelStatusCreate
	one.Link = updateRes.Link

	return one, nil
}
