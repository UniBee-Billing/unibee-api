package handler

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/gateway/ro"
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
		dao.Subscription.Columns().GmtModify:    gtime.Now(),
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
	Payment                   *entity.Payment                           `json:"payment" `
	ChannelSubscriptionDetail *ro.ChannelDetailSubscriptionInternalResp `json:"channelSubscriptionDetail"`
	ChannelInvoiceDetail      *ro.ChannelDetailInvoiceInternalResp      `json:"channelInvoiceDetail"`
	ChannelPaymentId          string                                    `json:"channelPaymentId" `
	ChannelSubscriptionId     string                                    `json:"channelSubscriptionId" `
	ChannelInvoiceId          string                                    `json:"channelInvoiceId"`
	ChannelUpdateId           string                                    `json:"channelUpdateId"`
	Status                    consts.SubscriptionStatusEnum             `json:"status"`
	ChannelStatus             string                                    `json:"channelStatus"                  `
	Data                      string                                    `json:"data"`
	ChannelItemData           string                                    `json:"channelItemData"`
	CancelAtPeriodEnd         bool                                      `json:"cancelAtPeriodEnd"`
	CurrentPeriodEnd          int64                                     `json:"currentPeriodEnd"`
	CurrentPeriodStart        int64                                     `json:"currentPeriodStart"`
	TrialEnd                  int64                                     `json:"trialEnd"`
}

func FinishPendingUpdateForSubscription(ctx context.Context, sub *entity.Subscription, channelPaymentId string, one *entity.SubscriptionPendingUpdate) (bool, error) {
	// 先创建 SubscriptionTimeLine 在做 Sub 更新
	err := CreateOrUpdateSubscriptionTimeline(ctx, sub, channelPaymentId)
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
	eiPendingSubUpdate := query.GetUnfinishedEffectImmediateSubscriptionPendingUpdateByChannelUpdateId(ctx, req.ChannelUpdateId)
	if eiPendingSubUpdate != nil {
		//更新单支付成功, EffectImmediate=true 需要用户 3DS 验证等场景
		_, err := FinishPendingUpdateForSubscription(ctx, sub, req.ChannelPaymentId, eiPendingSubUpdate)
		if err != nil {
			return err
		}
	} else {
		var byUpdate = false
		if len(sub.PendingUpdateId) > 0 {
			//有 pending 的更新单存在，检查支付是否对应更新单
			pendingSubUpdate := query.GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx, sub.PendingUpdateId)
			if pendingSubUpdate.UpdateAmount == req.Payment.PaymentFee {
				//金额一致
				_, err := FinishPendingUpdateForSubscription(ctx, sub, req.ChannelPaymentId, pendingSubUpdate)
				if err != nil {
					return err
				}
				byUpdate = true
			}
		}
		if !byUpdate {
			err := CreateOrUpdateSubscriptionTimeline(ctx, sub, req.ChannelPaymentId)
			if err != nil {
				g.Log().Errorf(ctx, "CreateOrUpdateSubscriptionTimeline error:%s", err.Error())
			}
		}
	}
	err := UpdateSubWithPayment(ctx, sub, req.ChannelSubscriptionDetail)
	if err != nil {
		return err
	}
	err = CreateOrUpdateInvoiceFromSubscriptionPaymentSuccess(ctx, sub.SubscriptionId, req.Payment, req.ChannelInvoiceDetail)
	if err != nil {
		return err
	}
	return nil
}

type SubscriptionPaymentFailureWebHookReq struct {
	Payment                   *entity.Payment                           `json:"payment" `
	ChannelSubscriptionDetail *ro.ChannelDetailSubscriptionInternalResp `json:"channelSubscriptionDetail"`
	ChannelInvoiceDetail      *ro.ChannelDetailInvoiceInternalResp      `json:"channelInvoiceDetail"`
	ChannelPaymentId          string                                    `json:"channelPaymentId" `
	ChannelSubscriptionId     string                                    `json:"channelSubscriptionId" `
	ChannelInvoiceId          string                                    `json:"channelInvoiceId"`
	ChannelUpdateId           string                                    `json:"channelUpdateId"`
}

func HandleSubscriptionPaymentFailure(ctx context.Context, req *SubscriptionPaymentFailureWebHookReq) error {
	sub := query.GetSubscriptionByChannelSubscriptionId(ctx, req.ChannelSubscriptionId)
	if sub == nil {
		return gerror.Newf("HandleSubscriptionPaymentFailure sub not found %s", req.ChannelSubscriptionId)
	}

	eiPendingSubUpdate := query.GetUnfinishedEffectImmediateSubscriptionPendingUpdateByChannelUpdateId(ctx, req.ChannelUpdateId)
	if eiPendingSubUpdate != nil {
		//更新单支付失败, EffectImmediate=true 需要用户 3DS 验证等场景

	} else {
		var byUpdate = false
		if len(sub.PendingUpdateId) > 0 {
			//有 pending 的更新单存在，检查支付是否对应更新单
			pendingSubUpdate := query.GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx, sub.PendingUpdateId)
			if pendingSubUpdate.UpdateAmount == req.Payment.PaymentFee {
				//金额一致

				byUpdate = true
			}
		}
		if !byUpdate {
			//没有匹配到更新单

		}
	}
	return nil
}
