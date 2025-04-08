package system

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	"strconv"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	event2 "unibee/internal/consumer/webhook/event"
	"unibee/internal/consumer/webhook/message"
	dao "unibee/internal/dao/default"
	detail2 "unibee/internal/logic/subscription/service/detail"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/system/subscription"
)

func (c *ControllerSubscription) InternalWebhookSync(ctx context.Context, req *subscription.InternalWebhookSyncReq) (res *subscription.InternalWebhookSyncRes, err error) {
	if req.IsSynchronous {
		total, firstId, lastId := syncSubscription(ctx, req)
		g.Log().Infof(ctx, "InternalWebhookSync Sync Finished with \nInternalWebhookSync req:%s \nInternalWebhookSync total:%d,firstId:%s,lastId:%s", utility.MarshalToJsonString(req), total, firstId, lastId)
		return &subscription.InternalWebhookSyncRes{Total: total, FirstId: firstId, LastId: lastId}, nil
	} else {
		go func() {
			backgroundCtx := context.Background()
			defer func() {
				if exception := recover(); exception != nil {
					if v, ok := exception.(error); ok && gerror.HasStack(v) {
						err = v
					} else {
						err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
					}
					g.Log().Errorf(backgroundCtx, "CreateOrUpdateInvoiceByChannelDetail Background Generate PDF panic error:%s\n", err.Error())
					return
				}
			}()
			total, firstId, lastId := syncSubscription(backgroundCtx, req)
			g.Log().Infof(backgroundCtx, "InternalWebhookSync Async Finished with \nInternalWebhookSync req:%s \nInternalWebhookSync total:%d,firstId:%s,lastId:%s", utility.MarshalToJsonString(req), total, firstId, lastId)
		}()
	}
	return &subscription.InternalWebhookSyncRes{}, nil
}

func syncSubscription(ctx context.Context, req *subscription.InternalWebhookSyncReq) (total int, firstId string, lastId string) {
	var count = 100
	var page = 0
	for {
		var list []*entity.Subscription
		query := dao.Subscription.Ctx(ctx)
		if req.StartId != nil {
			query = query.WhereGTE(dao.Subscription.Columns().Id, req.StartId)
		} else if req.StartTime != nil {
			query = query.WhereGTE(dao.Subscription.Columns().CreateTime, req.StartTime)
		}
		if req.EndId != nil {
			query = query.WhereLTE(dao.Subscription.Columns().Id, req.EndId)
		} else if req.EndTime != nil {
			query = query.WhereLTE(dao.Subscription.Columns().CreateTime, req.EndTime)
		}
		query = query.WhereIn(dao.Subscription.Columns().IsDeleted, []int{0}).
			Limit(page*count, count).
			OmitEmpty()
		_ = query.Scan(&list)
		if page == 0 && list != nil && len(list) > 0 {
			firstId = strconv.FormatUint(list[0].Id, 10)
		}
		if list != nil && len(list) > 0 {
			lastId = strconv.FormatUint(list[len(list)-1].Id, 10)
		}

		{
			for _, one := range list {
				event := event2.UNIBEE_WEBHOOK_EVENT_USER_CREATED
				if one.Status == consts.SubStatusActive || one.Status == consts.SubStatusIncomplete {
					event = event2.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_ACTIVATED
				} else if one.Status == consts.SubStatusCancelled {
					event = event2.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_CANCELLED
				} else if one.Status == consts.SubStatusExpired {
					event = event2.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_EXPIRED
				} else if one.Status == consts.SubStatusFailed {
					event = event2.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_FAILED
				}
				subDetailRes, err := detail2.SubscriptionDetail(ctx, one.SubscriptionId)
				if err != nil {
					continue
				}
				_, _ = redismq.Send(&redismq.Message{
					Topic: redismq2.TopicInternalWebhook.Topic,
					Tag:   redismq2.TopicInternalWebhook.Tag,
					Body: utility.MarshalToJsonString(&message.WebhookMessage{
						Id:         one.Id,
						Event:      event2.WebhookEvent(event),
						EventId:    utility.CreateEventId(),
						MerchantId: one.MerchantId,
						Data:       utility.FormatToGJson(subDetailRes),
					}),
				})
			}
		}
		total = total + len(list)
		// next page
		page = page + 1
		if list == nil || len(list) == 0 {
			break
		}
		g.Log().Infof(ctx, "InternalWebhookSync FinishedPage:%d with \nInternalWebhookSync req:%s \nInternalWebhookSync total:%d,firstId:%s,lastId:%s", page, utility.MarshalToJsonString(req), total, firstId, lastId)
	}
	return total, firstId, lastId
}
