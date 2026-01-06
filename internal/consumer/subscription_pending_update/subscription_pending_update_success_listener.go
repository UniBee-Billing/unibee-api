package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	subscription "unibee/internal/consumer/webhook/subscription_pending_update"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/invoice/service"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type SubscriptionPendingUpdateSuccessListener struct {
}

func (t SubscriptionPendingUpdateSuccessListener) GetTopic() string {
	return redismq2.TopicSubscriptionPendingUpdateSuccess.Topic
}

func (t SubscriptionPendingUpdateSuccessListener) GetTag() string {
	return redismq2.TopicSubscriptionPendingUpdateSuccess.Tag
}

func (t SubscriptionPendingUpdateSuccessListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Infof(ctx, "SubscriptionPendingUpdateSuccessListener Receive Message:%s", utility.MarshalToJsonString(message))
	pendingUpdate := query.GetSubscriptionPendingUpdateByPendingUpdateId(ctx, message.Body)
	if pendingUpdate != nil {
		subscription.SendMerchantSubscriptionPendingUpdateWebhookBackground(pendingUpdate, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_PENDING_UPDATE_SUCCESS, message.CustomData)
		if len(pendingUpdate.SubscriptionId) > 0 {
			var list []*entity.Invoice
			_ = dao.Invoice.Ctx(ctx).
				Where(dao.Invoice.Columns().SubscriptionId, pendingUpdate.SubscriptionId).
				Where(dao.Invoice.Columns().BizType, consts.BizTypeSubscription).
				OrderDesc(dao.Invoice.Columns().Id).
				Limit(10).
				OmitEmpty().Scan(&list)
			for _, one := range list {
				if one != nil && one.Status == consts.InvoiceStatusProcessing && one.CreateFrom == consts.InvoiceAutoChargeFlag {
					err := service.CancelProcessingInvoice(ctx, one.InvoiceId, "TryCancelSubscriptionLatestAutoChargeInvoice")
					if err != nil {
						g.Log().Errorf(ctx, `SubscriptionPendingUpdateSuccessListener TryCancelSubscriptionLatestAutoChargeInvoice failure error:%s`, err.Error())
					}
				}
			}
		}
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionPendingUpdateSuccessListener())
	fmt.Println("SubscriptionPendingUpdateSuccessListener RegisterListener")
}

func NewSubscriptionPendingUpdateSuccessListener() *SubscriptionPendingUpdateSuccessListener {
	return &SubscriptionPendingUpdateSuccessListener{}
}
