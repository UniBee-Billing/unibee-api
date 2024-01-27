package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/channel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
)

func HandleSubscriptionWebhookEvent(ctx context.Context, subscription *entity.Subscription, eventType string, details *ro.ChannelDetailSubscriptionInternalResp) error {
	//更新 Subscription
	return UpdateSubWithChannelDetailBack(ctx, subscription, details)
}

func UpdateSubWithChannelDetailBack(ctx context.Context, subscription *entity.Subscription, details *ro.ChannelDetailSubscriptionInternalResp) error {
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
		dao.Subscription.Columns().Status:                 details.Status,
		dao.Subscription.Columns().ChannelSubscriptionId:  details.ChannelSubscriptionId,
		dao.Subscription.Columns().ChannelStatus:          details.ChannelStatus,
		dao.Subscription.Columns().ChannelItemData:        details.ChannelItemData,
		dao.Subscription.Columns().CancelAtPeriodEnd:      cancelAtPeriodEnd,
		dao.Subscription.Columns().BillingCycleAnchor:     details.BillingCycleAnchor,
		dao.Subscription.Columns().ChannelLatestInvoiceId: details.ChannelLatestInvoiceId,
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

func UpdateSubWithPayment(ctx context.Context, subscription *entity.Subscription, details *ro.ChannelDetailSubscriptionInternalResp) error {
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

func FinishPendingUpdateForSubscription(ctx context.Context, sub *entity.Subscription, one *entity.SubscriptionPendingUpdate) (bool, error) {
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

func HandleSubscriptionPaymentSuccess(ctx context.Context, req *SubscriptionPaymentSuccessWebHookReq) error {
	sub := query.GetSubscriptionByChannelSubscriptionId(ctx, req.ChannelSubscriptionId)
	if sub == nil {
		return gerror.Newf("HandleSubscriptionPaymentSuccess sub not found %s", req.ChannelSubscriptionId)
	}
	eiPendingSubUpdate := query.GetUnfinishedEffectImmediateSubscriptionPendingUpdateByChannelUpdateId(ctx, req.ChannelSubscriptionUpdateId)
	if eiPendingSubUpdate != nil {
		//更新单支付成功, EffectImmediate=true 需要用户 3DS 验证等场景
		_, err := FinishPendingUpdateForSubscription(ctx, sub, eiPendingSubUpdate)
		if err != nil {
			return err
		}
	} else {
		var byUpdate = false
		if len(sub.PendingUpdateId) > 0 {
			//有 pending 的更新单存在，检查支付是否对应更新单
			pendingSubUpdate := query.GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx, sub.PendingUpdateId)
			if pendingSubUpdate.UpdateAmount == req.Payment.TotalAmount {
				//金额一致
				_, err := FinishPendingUpdateForSubscription(ctx, sub, pendingSubUpdate)
				if err != nil {
					return err
				}
				byUpdate = true
			}
		}
		if !byUpdate && req.Payment.TotalAmount == sub.Amount {
			// billing-cycle
			err := CreateOrUpdateSubscriptionTimeline(ctx, sub, fmt.Sprintf("cycle-paymentId-%s", req.Payment.PaymentId))
			if err != nil {
				g.Log().Errorf(ctx, "CreateOrUpdateSubscriptionTimeline error:%s", err.Error())
			}
		}
	}
	err := UpdateSubWithPayment(ctx, sub, req.ChannelSubscriptionDetail)
	if err != nil {
		return err
	}
	//重新获取
	sub = query.GetSubscriptionByChannelSubscriptionId(ctx, req.ChannelSubscriptionId)
	// todo generate InvoiceSimplify and update payment

	//err = handler.CreateOrUpdateInvoiceForSubscriptionPaymentSuccess(ctx, &handler.CreateInvoiceInternalReq{
	//	Payment:                          req.Payment,
	//	ChannelInvoiceId:                 req.ChannelInvoiceId,
	//	Currency:                         sub.Currency,
	//	PlanId:                           sub.PlanId,
	//	Quantity:                         sub.Quantity,
	//	AddonJsonData:                    sub.AddonData,
	//	TaxScale:                         sub.TaxScale,
	//	UserId:                           sub.UserId,
	//	MerchantId:                       sub.MerchantId,
	//	SubscriptionId:                   sub.SubscriptionId,
	//	ChannelId:                        sub.ChannelId,
	//	InvoiceStatus:                    consts.InvoiceStatusPaid,
	//	ChannelDetailInvoiceInternalResp: req.ChannelInvoiceDetail,
	//	PeriodStart:                      sub.CurrentPeriodStart,
	//	PeriodEnd:                        sub.CurrentPeriodEnd,
	//})
	if err != nil {
		return err
	}
	return nil
}

type SubscriptionPaymentFailureWebHookReq struct {
	Payment                     *entity.Payment                           `json:"payment" `
	ChannelSubscriptionDetail   *ro.ChannelDetailSubscriptionInternalResp `json:"channelSubscriptionDetail"`
	ChannelInvoiceDetail        *ro.ChannelDetailInvoiceInternalResp      `json:"channelInvoiceDetail"`
	ChannelPaymentId            string                                    `json:"channelPaymentId" `
	ChannelSubscriptionId       string                                    `json:"channelSubscriptionId" `
	ChannelInvoiceId            string                                    `json:"channelInvoiceId"`
	ChannelSubscriptionUpdateId string                                    `json:"channelUpdateId"`
}

func HandleSubscriptionPaymentFailure(ctx context.Context, req *SubscriptionPaymentFailureWebHookReq) error {
	sub := query.GetSubscriptionByChannelSubscriptionId(ctx, req.ChannelSubscriptionId)
	if sub == nil {
		return gerror.Newf("HandleSubscriptionPaymentFailure sub not found %s", req.ChannelSubscriptionId)
	}
	return nil
}

func HandleSubscriptionPaymentWaitAuthorized(ctx context.Context, req *SubscriptionPaymentFailureWebHookReq) error {
	sub := query.GetSubscriptionByChannelSubscriptionId(ctx, req.ChannelSubscriptionId)
	if sub == nil {
		return gerror.Newf("HandleSubscriptionPaymentWaitAuthorized sub not found %s", req.ChannelSubscriptionId)
	}

	//eiPendingSubUpdate := query.GetUnfinishedEffectImmediateSubscriptionPendingUpdateByChannelUpdateId(ctx, req.ChannelSubscriptionUpdateId)
	//if eiPendingSubUpdate != nil {
	//	//更新单支付失败, EffectImmediate=true 需要用户 3DS 验证等场景
	//	//金额一致
	//	err := handler.CreateOrUpdateInvoiceForSubscriptionPaymentSuccess(ctx, &handler.CreateInvoiceInternalReq{
	//		Payment:                          req.Payment,
	//		Currency:                         eiPendingSubUpdate.UpdateCurrency,
	//		PlanId:                           eiPendingSubUpdate.UpdatePlanId,
	//		Quantity:                         eiPendingSubUpdate.UpdateQuantity,
	//		AddonJsonData:                    eiPendingSubUpdate.UpdateAddonData,
	//		TaxScale:                         sub.TaxScale,
	//		UserId:                           sub.UserId,
	//		MerchantId:                       sub.MerchantId,
	//		SubscriptionId:                   sub.SubscriptionId,
	//		ChannelId:                        sub.ChannelId,
	//		InvoiceStatus:                    consts.InvoiceStatusProcessing,
	//		ChannelDetailInvoiceInternalResp: req.ChannelInvoiceDetail,
	//		PeriodStart:                      sub.CurrentPeriodStart, // todo mark 周期不确定
	//		PeriodEnd:                        sub.CurrentPeriodEnd,
	//	})
	//	if err != nil {
	//		return err
	//	}
	//} else {
	//	var byUpdate = false
	//	if len(sub.PendingUpdateId) > 0 {
	//		//有 pending 的更新单存在，检查支付是否对应更新单
	//		pendingSubUpdate := query.GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx, sub.PendingUpdateId)
	//		if pendingSubUpdate.UpdateAmount == req.Payment.TotalAmount {
	//			//金额一致
	//			err := handler.CreateOrUpdateInvoiceForSubscriptionPaymentSuccess(ctx, &handler.CreateInvoiceInternalReq{
	//				Payment:                          req.Payment,
	//				Currency:                         pendingSubUpdate.UpdateCurrency,
	//				PlanId:                           pendingSubUpdate.UpdatePlanId,
	//				Quantity:                         pendingSubUpdate.UpdateQuantity,
	//				AddonJsonData:                    pendingSubUpdate.UpdateAddonData,
	//				TaxScale:                         sub.TaxScale,
	//				UserId:                           sub.UserId,
	//				MerchantId:                       sub.MerchantId,
	//				SubscriptionId:                   sub.SubscriptionId,
	//				ChannelId:                        sub.ChannelId,
	//				InvoiceStatus:                    consts.InvoiceStatusProcessing,
	//				ChannelDetailInvoiceInternalResp: req.ChannelInvoiceDetail,
	//				PeriodStart:                      sub.CurrentPeriodStart, // todo mark 周期不确定
	//				PeriodEnd:                        sub.CurrentPeriodEnd,
	//			})
	//			if err != nil {
	//				return err
	//			}
	//			byUpdate = true
	//		}
	//	}
	//	if !byUpdate {
	//		//没有匹配到更新单
	//		err := handler.CreateOrUpdateInvoiceForSubscriptionPaymentSuccess(ctx, &handler.CreateInvoiceInternalReq{
	//			Payment:                          req.Payment,
	//			Currency:                         sub.Currency,
	//			PlanId:                           sub.PlanId,
	//			Quantity:                         sub.Quantity,
	//			AddonJsonData:                    sub.AddonData,
	//			TaxScale:                         sub.TaxScale,
	//			UserId:                           sub.UserId,
	//			MerchantId:                       sub.MerchantId,
	//			SubscriptionId:                   sub.SubscriptionId,
	//			ChannelId:                        sub.ChannelId,
	//			InvoiceStatus:                    consts.InvoiceStatusProcessing,
	//			ChannelDetailInvoiceInternalResp: req.ChannelInvoiceDetail,
	//			PeriodStart:                      sub.CurrentPeriodEnd,
	//			PeriodEnd:                        sub.CurrentPeriodEnd + (sub.CurrentPeriodEnd - sub.CurrentPeriodStart), // + 1 周期 todo mark 确认
	//		})
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}
	return nil
}
