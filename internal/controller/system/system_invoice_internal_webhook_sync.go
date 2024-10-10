package system

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	"unibee/api/bean/detail"
	"unibee/api/system/invoice"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	event2 "unibee/internal/consumer/webhook/event"
	"unibee/internal/consumer/webhook/message"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func (c *ControllerInvoice) InternalWebhookSync(ctx context.Context, req *invoice.InternalWebhookSyncReq) (res *invoice.InternalWebhookSyncRes, err error) {
	if req.IsSynchronous {
		total, firstId, lastId := syncInvoice(ctx, req)
		g.Log().Infof(ctx, "InternalWebhookSync Sync Finished with \nInternalWebhookSync req:%s \nInternalWebhookSync total:%d,firstId:%s,lastId:%s", utility.MarshalToJsonString(req), total, firstId, lastId)
		return &invoice.InternalWebhookSyncRes{Total: total, FirstId: firstId, LastId: lastId}, nil
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
			total, firstId, lastId := syncInvoice(backgroundCtx, req)
			g.Log().Infof(backgroundCtx, "InternalWebhookSync Async Finished with \nInternalWebhookSync req:%s \nInternalWebhookSync total:%d,firstId:%s,lastId:%s", utility.MarshalToJsonString(req), total, firstId, lastId)
		}()
	}
	return &invoice.InternalWebhookSyncRes{}, nil
}

func syncInvoice(ctx context.Context, req *invoice.InternalWebhookSyncReq) (total int, firstId string, lastId string) {
	var count = 100
	var page = 0
	for {
		var list []*entity.Invoice
		query := dao.Invoice.Ctx(ctx)
		if req.StartId != nil {
			query = query.WhereGTE(dao.Invoice.Columns().Id, req.StartId)
		} else if req.StartTime != nil {
			query = query.WhereGTE(dao.Invoice.Columns().CreateTime, req.StartTime)
		}
		if req.EndId != nil {
			query = query.WhereLTE(dao.Invoice.Columns().Id, req.EndId)
		} else if req.EndTime != nil {
			query = query.WhereLTE(dao.Invoice.Columns().CreateTime, req.EndTime)
		}
		query = query.WhereIn(dao.Invoice.Columns().IsDeleted, []int{0}).
			Limit(page*count, count).
			OmitEmpty()
		_ = query.Scan(&list)
		if page == 0 && list != nil && len(list) > 0 {
			firstId = list[0].InvoiceId
		}
		if list != nil && len(list) > 0 {
			lastId = list[len(list)-1].InvoiceId
		}

		{
			for _, one := range list {
				event := event2.UNIBEE_WEBHOOK_EVENT_INVOICE_CREATED
				if one.Status == consts.InvoiceStatusPaid {
					event = event2.UNIBEE_WEBHOOK_EVENT_INVOICE_PAID
				} else if one.Status == consts.InvoiceStatusCancelled {
					event = event2.UNIBEE_WEBHOOK_EVENT_INVOICE_CANCELLED
				} else if one.Status == consts.InvoiceStatusFailed {
					event = event2.UNIBEE_WEBHOOK_EVENT_INVOICE_FAILED
				} else if one.Status == consts.InvoiceStatusProcessing {
					event = event2.UNIBEE_WEBHOOK_EVENT_INVOICE_PROCESS
				}
				_, _ = redismq.Send(&redismq.Message{
					Topic: redismq2.TopicInternalWebhook.Topic,
					Tag:   redismq2.TopicInternalWebhook.Tag,
					Body: utility.MarshalToJsonString(&message.WebhookMessage{
						Id:         one.Id,
						Event:      event2.WebhookEvent(event),
						EventId:    utility.CreateEventId(),
						MerchantId: one.MerchantId,
						Data:       utility.FormatToGJson(detail.ConvertInvoiceToDetail(ctx, one)),
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
	}
	return total, firstId, lastId
}
