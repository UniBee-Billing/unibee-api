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
	update, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                 details.Status,
		dao.Subscription.Columns().ChannelSubscriptionId:  details.ChannelSubscriptionId,
		dao.Subscription.Columns().ChannelStatus:          details.ChannelStatus,
		dao.Subscription.Columns().ChannelLatestInvoiceId: details.ChannelLatestInvoiceId,
		dao.Subscription.Columns().ChannelItemData:        details.ChannelItemData,
		dao.Subscription.Columns().CancelAtPeriodEnd:      cancelAtPeriodEnd,
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
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return gerror.Newf("HandleSubscriptionWebhookEvent err:%s", update)
	}
	//处理更新事件 todo mark

	return nil
}

type SubscriptionPaymentSuccessWebHookReq struct {
	Payment               *entity.Payment               `json:"payment" `
	ChannelPaymentId      string                        `json:"channelPaymentId" `
	ChannelSubscriptionId string                        `json:"channelSubscriptionId" `
	ChannelInvoiceId      string                        `json:"channelInvoiceId"`
	ChannelUpdateId       string                        `json:"channelUpdateId"`
	Status                consts.SubscriptionStatusEnum `json:"status"`
	ChannelStatus         string                        `json:"channelStatus"                  ` // 货币
	Data                  string                        `json:"data"`
	ChannelItemData       string                        `json:"channelItemData"`
	CancelAtPeriodEnd     bool                          `json:"cancelAtPeriodEnd"`
	CurrentPeriodEnd      int64                         `json:"currentPeriodEnd"`
	CurrentPeriodStart    int64                         `json:"currentPeriodStart"`
	TrialEnd              int64                         `json:"trialEnd"`
}

func FinishPendingUpdateForSubscription(ctx context.Context, one *entity.SubscriptionPendingUpdate) (bool, error) {
	// todo 使用事务
	update, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().PlanId:    one.UpdatePlanId,
		dao.Subscription.Columns().Quantity:  one.UpdateQuantity,
		dao.Subscription.Columns().AddonData: one.UpdateAddonData,
		dao.Subscription.Columns().Amount:    one.UpdateAmount,
		dao.Subscription.Columns().Currency:  one.UpdateCurrency,
		dao.Subscription.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, one.SubscriptionId).OmitNil().Update()
	if err != nil {
		return false, err
	}
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return false, gerror.Newf("SubscriptionPendingUpdate update subscription err:%s", update)
	}
	_, err = dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPendingUpdate.Columns().Status:    consts.PendingSubStatusFinished,
		dao.SubscriptionPendingUpdate.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionPendingUpdate.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return false, err
	}
	// todo mark sub 更新成功，生成 Proration Invoice 发票、 timeline等后续流程
	return true, nil
}

func HandleSubscriptionPaymentSuccess(ctx context.Context, req *SubscriptionPaymentSuccessWebHookReq) error {
	//sub := query.GetSubscriptionByChannelSubscriptionId(ctx, req.ChannelSubscriptionId)
	pendingSubUpdate := query.GetCreatedSubscriptionPendingUpdateByChannelUpdateId(ctx, req.ChannelUpdateId)
	// todo mark 处理逻辑需要优化，降级或者不马上生效的更新，是在下一周期生效
	if pendingSubUpdate != nil {
		//更新单支付成功
		//todo mark 更新之前 先 校验渠道 Sub Plan 是否已更改， 金额是否一致
		_, err := FinishPendingUpdateForSubscription(ctx, pendingSubUpdate)
		if err != nil {
			return err
		}
		//更新 SubscriptionPendingUpdate
		_, err = dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
			dao.SubscriptionPendingUpdate.Columns().Status:    consts.PendingSubStatusInit,
			dao.SubscriptionPendingUpdate.Columns().GmtModify: gtime.Now(),
		}).Where(dao.SubscriptionPendingUpdate.Columns().Id, pendingSubUpdate.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		sub := query.GetSubscriptionByChannelSubscriptionId(ctx, req.ChannelSubscriptionId)
		if sub != nil {
			one := &entity.SubscriptionTimeline{
				MerchantId:      sub.MerchantId,
				UserId:          sub.UserId,
				SubscriptionId:  sub.SubscriptionId,
				InvoiceId:       "", // todo mark
				UniqueId:        req.Payment.PaymentId,
				Currency:        sub.Currency,
				PlanId:          sub.PlanId,
				Quantity:        sub.Quantity,
				AddonData:       sub.AddonData,
				ChannelId:       sub.ChannelId,
				PeriodStart:     gtime.Now().Timestamp(),
				PeriodEnd:       sub.CurrentPeriodEnd,
				PeriodStartTime: gtime.Now(),
				PeriodEndTime:   sub.CurrentPeriodEndTime,
				PaymentId:       req.Payment.PaymentId,
			}

			_, err = dao.SubscriptionTimeline.Ctx(ctx).Data(one).OmitNil().Insert(one)
			if err != nil {
				err = gerror.Newf(`HandleSubscriptionPaymentSuccess record insert failure %s`, err.Error())
				return err
			}
		}
	} else {
		sub := query.GetSubscriptionByChannelSubscriptionId(ctx, req.ChannelSubscriptionId)
		if sub != nil {
			// todo mark sub Next Period，生成 对应 Invoice 发票、 timeline等后续流程
			one := &entity.SubscriptionTimeline{
				MerchantId:      sub.MerchantId,
				UserId:          sub.UserId,
				SubscriptionId:  sub.SubscriptionId,
				InvoiceId:       "", // todo mark
				UniqueId:        req.Payment.PaymentId,
				Currency:        sub.Currency,
				PlanId:          sub.PlanId,
				Quantity:        sub.Quantity,
				AddonData:       sub.AddonData,
				ChannelId:       sub.ChannelId,
				PeriodStart:     sub.CurrentPeriodStart,
				PeriodEnd:       sub.CurrentPeriodEnd,
				PeriodStartTime: sub.CurrentPeriodStartTime,
				PeriodEndTime:   sub.CurrentPeriodEndTime,
				PaymentId:       req.Payment.PaymentId,
			}

			_, err := dao.SubscriptionTimeline.Ctx(ctx).Data(one).OmitNil().Insert(one)
			if err != nil {
				err = gerror.Newf(`HandleSubscriptionPaymentSuccess record insert failure %s`, err.Error())
				return err
			}
		}
	}
	return nil
}
