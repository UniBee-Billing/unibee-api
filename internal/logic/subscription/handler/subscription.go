package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/channel/out"
	"go-oversea-pay/internal/logic/channel/ro"
	"go-oversea-pay/internal/logic/invoice/handler"
	"go-oversea-pay/internal/logic/invoice/invoice_compute"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func HandleSubscriptionWebhookEvent(ctx context.Context, subscription *entity.Subscription, eventType string, details *ro.ChannelDetailSubscriptionInternalResp) error {
	//更新 Subscription
	return UpdateSubWithChannelDetailBack(ctx, subscription, details)
}

func UpdateSubWithChannelDetailBack(ctx context.Context, subscription *entity.Subscription, details *ro.ChannelDetailSubscriptionInternalResp) error {
	if consts.SubscriptionCycleUnderUniBeeControl {
		// not sync attribute from channel
		return nil
	}
	var cancelAtPeriodEnd = 0
	if details.CancelAtPeriodEnd {
		cancelAtPeriodEnd = 1
	}
	var firstPayTime *gtime.Time
	if subscription.FirstPayTime == nil && details.Status == consts.SubStatusActive {
		firstPayTime = gtime.Now()
	}
	var gmtModify = subscription.GmtModify
	if subscription.Status != int(details.Status) ||
		subscription.CancelAtPeriodEnd != cancelAtPeriodEnd ||
		subscription.BillingCycleAnchor != details.BillingCycleAnchor ||
		subscription.TrialEnd != details.TrialEnd ||
		subscription.FirstPayTime != firstPayTime {
		gmtModify = gtime.Now()
	}

	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                      details.Status,
		dao.Subscription.Columns().ChannelSubscriptionId:       details.ChannelSubscriptionId,
		dao.Subscription.Columns().ChannelStatus:               details.ChannelStatus,
		dao.Subscription.Columns().ChannelItemData:             details.ChannelItemData,
		dao.Subscription.Columns().CancelAtPeriodEnd:           cancelAtPeriodEnd,
		dao.Subscription.Columns().BillingCycleAnchor:          details.BillingCycleAnchor,
		dao.Subscription.Columns().ChannelLatestInvoiceId:      details.ChannelLatestInvoiceId,
		dao.Subscription.Columns().ChannelDefaultPaymentMethod: details.ChannelDefaultPaymentMethod,
		//dao.Subscription.Columns().CurrentPeriodStart:     details.CurrentPeriodStart,
		//dao.Subscription.Columns().CurrentPeriodEnd:       details.CurrentPeriodEnd,
		//dao.Subscription.Columns().CurrentPeriodStartTime: gtime.NewFromTimeStamp(details.CurrentPeriodStart),
		//dao.Subscription.Columns().CurrentPeriodEndTime:   gtime.NewFromTimeStamp(details.CurrentPeriodEnd),
		dao.Subscription.Columns().TrialEnd:     details.TrialEnd,
		dao.Subscription.Columns().GmtModify:    gmtModify,
		dao.Subscription.Columns().FirstPayTime: firstPayTime,
	}).Where(dao.Subscription.Columns().Id, subscription.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	//rowAffected, err := update.RowsAffected()
	//if rowAffected != 1 {
	//	return gerror.Newf("HandleSubscriptionWebhookEvent err:%s", update)
	//}
	//处理更新事件 todo mark

	return nil
}

func UpdateSubscriptionBillingCycleWithPayment(ctx context.Context, subscription *entity.Subscription, details *ro.ChannelDetailSubscriptionInternalResp) error {
	var cancelAtPeriodEnd = 0
	if details.CancelAtPeriodEnd {
		cancelAtPeriodEnd = 1
	}
	var firstPayTime *gtime.Time
	if subscription.FirstPayTime == nil && details.Status == consts.SubStatusActive {
		firstPayTime = gtime.Now()
	}
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                 details.Status,
		dao.Subscription.Columns().ChannelSubscriptionId:  details.ChannelSubscriptionId,
		dao.Subscription.Columns().ChannelStatus:          details.ChannelStatus,
		dao.Subscription.Columns().ChannelItemData:        details.ChannelItemData,
		dao.Subscription.Columns().CancelAtPeriodEnd:      cancelAtPeriodEnd,
		dao.Subscription.Columns().BillingCycleAnchor:     details.BillingCycleAnchor,
		dao.Subscription.Columns().ChannelLatestInvoiceId: details.ChannelLatestInvoiceId,
		dao.Subscription.Columns().CurrentPeriodStart:     details.CurrentPeriodStart,
		dao.Subscription.Columns().CurrentPeriodEnd:       details.CurrentPeriodEnd,
		dao.Subscription.Columns().CurrentPeriodStartTime: gtime.NewFromTimeStamp(details.CurrentPeriodStart),
		dao.Subscription.Columns().CurrentPeriodEndTime:   gtime.NewFromTimeStamp(details.CurrentPeriodEnd),
		dao.Subscription.Columns().TrialEnd:               details.TrialEnd,
		dao.Subscription.Columns().GmtModify:              gtime.Now(),
		dao.Subscription.Columns().FirstPayTime:           firstPayTime,
	}).Where(dao.Subscription.Columns().Id, subscription.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	//处理更新事件 todo mark

	return nil
}

type SubscriptionPaymentSuccessWebHookReq struct {
	Payment                     *entity.Payment                           `json:"payment" `
	ChannelSubscriptionDetail   *ro.ChannelDetailSubscriptionInternalResp `json:"channelSubscriptionDetail"`
	ChannelInvoiceDetail        *ro.ChannelDetailInvoiceInternalResp      `json:"channelInvoiceDetail"`
	ChannelPaymentId            string                                    `json:"channelPaymentId" `
	ChannelInvoiceId            string                                    `json:"channelInvoiceId"`
	ChannelSubscriptionId       string                                    `json:"channelSubscriptionId" `
	ChannelSubscriptionUpdateId string                                    `json:"channelSubscriptionUpdateId"`
	Status                      consts.SubscriptionStatusEnum             `json:"status"`
	ChannelStatus               string                                    `json:"channelStatus"                  `
	Data                        string                                    `json:"data"`
	ChannelItemData             string                                    `json:"channelItemData"`
	CancelAtPeriodEnd           bool                                      `json:"cancelAtPeriodEnd"`
	CurrentPeriodEnd            int64                                     `json:"currentPeriodEnd"`
	CurrentPeriodStart          int64                                     `json:"currentPeriodStart"`
	TrialEnd                    int64                                     `json:"trialEnd"`
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

func FinishPendingUpdateForSubscription(ctx context.Context, sub *entity.Subscription, one *entity.SubscriptionPendingUpdate) (bool, error) {
	if one.Status == consts.PendingSubStatusFinished {
		return true, nil
	}
	if consts.ProrationUsingUniBeeCompute && one.EffectImmediate == 1 {
		var addonParams []*ro.SubscriptionPlanAddonParamRo
		err := utility.UnmarshalFromJsonString(one.UpdateAddonData, &addonParams)
		if err != nil {
			return false, err
		}
		_, err = out.GetPayChannelServiceProvider(ctx, one.ChannelId).DoRemoteChannelSubscriptionUpdate(ctx, &ro.ChannelUpdateSubscriptionInternalReq{
			Plan:            query.GetPlanById(ctx, one.UpdatePlanId),
			Quantity:        one.UpdateQuantity,
			AddonPlans:      checkAndListAddonsFromParams(ctx, addonParams, one.ChannelId),
			PlanChannel:     query.GetPlanChannel(ctx, one.UpdatePlanId, one.ChannelId),
			Subscription:    query.GetSubscriptionBySubscriptionId(ctx, one.SubscriptionId),
			ProrationDate:   one.ProrationDate,
			EffectImmediate: false,
		})
		if err != nil {
			return false, err
		}
	}

	// 先创建 SubscriptionTimeLine 在做 Sub 更新
	err := CreateOrUpdateSubscriptionTimeline(ctx, sub, fmt.Sprintf("pendingUpdateFinish-%s", one.UpdateSubscriptionId))
	if err != nil {
		g.Log().Errorf(ctx, "CreateOrUpdateSubscriptionTimeline error:%s", err.Error())
	}
	// todo mark 使用事务
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().PlanId:          one.UpdatePlanId,
		dao.Subscription.Columns().Quantity:        one.UpdateQuantity,
		dao.Subscription.Columns().AddonData:       one.UpdateAddonData,
		dao.Subscription.Columns().Amount:          one.UpdateAmount,
		dao.Subscription.Columns().Currency:        one.UpdateCurrency,
		dao.Subscription.Columns().LastUpdateTime:  gtime.Now().Timestamp(),
		dao.Subscription.Columns().GmtModify:       gtime.Now(),
		dao.Subscription.Columns().PendingUpdateId: "", //清除标记的更新单
	}).Where(dao.Subscription.Columns().SubscriptionId, one.SubscriptionId).OmitNil().Update()
	if err != nil {
		return false, err
	}
	_, err = dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPendingUpdate.Columns().Status:    consts.PendingSubStatusFinished,
		dao.SubscriptionPendingUpdate.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionPendingUpdate.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return false, err
	}

	return true, nil
}

func FinishNextBillingCycleForSubscription(ctx context.Context, sub *entity.Subscription, payment *entity.Payment, channelSubscriptionDetail *ro.ChannelDetailSubscriptionInternalResp) error {
	// billing-cycle
	err := CreateOrUpdateSubscriptionTimeline(ctx, sub, fmt.Sprintf("cycle-paymentId-%s", payment.PaymentId))
	if err != nil {
		g.Log().Errorf(ctx, "FinishNextBillingCycleForSubscription error:%s", err.Error())
	}
	err = UpdateSubscriptionBillingCycleWithPayment(ctx, sub, channelSubscriptionDetail)
	return err
}

func HandleSubscriptionPaymentUpdate(ctx context.Context, req *SubscriptionPaymentSuccessWebHookReq) error {
	sub := query.GetSubscriptionByChannelSubscriptionId(ctx, req.ChannelSubscriptionId)
	if sub == nil {
		return gerror.Newf("HandleSubscriptionPaymentUpdate sub not found %s", req.ChannelSubscriptionId)
	}
	eiPendingSubUpdate := query.GetUnfinishedEffectImmediateSubscriptionPendingUpdateByChannelUpdateId(ctx, req.ChannelSubscriptionUpdateId)
	if eiPendingSubUpdate != nil {
		// subscription_update
		//更新单支付成功, EffectImmediate=true 需要用户 3DS 验证等场景
		if req.Payment.Status == consts.PAY_SUCCESS {
			_, err := FinishPendingUpdateForSubscription(ctx, sub, eiPendingSubUpdate)
			if err != nil {
				return err
			}
		}
	} else {
		// subscription_cycle
		var byUpdate = false
		if len(sub.PendingUpdateId) > 0 {
			//有 pending 的更新单存在，检查支付是否对应更新单
			pendingSubUpdate := query.GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx, sub.PendingUpdateId)
			if pendingSubUpdate.UpdateAmount == req.Payment.TotalAmount {
				// compensate next pending update billing cycle invoice
				one := query.GetInvoiceByPaymentId(ctx, req.Payment.PaymentId)
				if one == nil {
					sub = query.GetSubscriptionByChannelSubscriptionId(ctx, req.ChannelSubscriptionId)

					var nextPeriodStart = sub.CurrentPeriodEnd
					if sub.TrialEnd > sub.CurrentPeriodEnd {
						nextPeriodStart = sub.TrialEnd
					}
					var nextPeriodEnd = nextPeriodStart + (sub.CurrentPeriodEnd - sub.CurrentPeriodStart)

					invoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
						Currency:      pendingSubUpdate.UpdateCurrency,
						PlanId:        pendingSubUpdate.UpdateAmount,
						Quantity:      pendingSubUpdate.UpdateQuantity,
						AddonJsonData: pendingSubUpdate.UpdateAddonData,
						TaxScale:      sub.TaxScale,
						PeriodStart:   nextPeriodStart,
						PeriodEnd:     nextPeriodEnd,
					})
					_ = handler.CreateOrUpdateInvoiceFromPayment(ctx, invoice, req.Payment, req.ChannelInvoiceDetail)
				}

				if req.Payment.Status == consts.PAY_SUCCESS {
					_, err := FinishPendingUpdateForSubscription(ctx, sub, pendingSubUpdate)
					if err != nil {
						return err
					}

					err = FinishNextBillingCycleForSubscription(ctx, sub, req.Payment, req.ChannelSubscriptionDetail)
					if err != nil {
						return err
					}
				}
				byUpdate = true
			}
		}
		if !byUpdate && req.Payment.TotalAmount == sub.Amount {
			// compensate billing cycle invoice
			one := query.GetInvoiceByPaymentId(ctx, req.Payment.PaymentId)
			if one == nil {
				sub = query.GetSubscriptionByChannelSubscriptionId(ctx, req.ChannelSubscriptionId)

				var nextPeriodStart = sub.CurrentPeriodEnd
				if sub.TrialEnd > sub.CurrentPeriodEnd {
					nextPeriodStart = sub.TrialEnd
				}
				var nextPeriodEnd = nextPeriodStart + (sub.CurrentPeriodEnd - sub.CurrentPeriodStart)

				invoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
					Currency:      sub.Currency,
					PlanId:        sub.PlanId,
					Quantity:      sub.Quantity,
					AddonJsonData: sub.AddonData,
					TaxScale:      sub.TaxScale,
					PeriodStart:   nextPeriodStart,
					PeriodEnd:     nextPeriodEnd,
				})
				_ = handler.CreateOrUpdateInvoiceFromPayment(ctx, invoice, req.Payment, req.ChannelInvoiceDetail)
			}
			if req.Payment.Status == consts.PAY_SUCCESS {
				err := FinishNextBillingCycleForSubscription(ctx, sub, req.Payment, req.ChannelSubscriptionDetail)
				if err != nil {
					return err
				}
			}
		} else if !byUpdate {
			// other situation payment need implementation todo mark

		}
	}
	return nil
}
