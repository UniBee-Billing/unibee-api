package handler

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
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
	update, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                 details.Status,
		dao.Subscription.Columns().ChannelStatus:          details.ChannelStatus,
		dao.Subscription.Columns().ChannelLatestInvoiceId: details.ChannelLatestInvoiceId,
		dao.Subscription.Columns().CancelAtPeriodEnd:      cancelAtPeriodEnd,
		dao.Subscription.Columns().GmtModify:              gtime.Now(),
	}).Where(dao.Subscription.Columns().Id, subscription.Id).OmitEmpty().Update()
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
