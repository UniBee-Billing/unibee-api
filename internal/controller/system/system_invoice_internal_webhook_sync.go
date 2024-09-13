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
	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				g.Log().Errorf(context.Background(), "CreateOrUpdateInvoiceByChannelDetail Background Generate PDF panic error:%s\n", err.Error())
				return
			}
		}()
		var list []*entity.Invoice
		var count = 100
		var page = 0
		backgroundCtx := context.Background()
		for {
			query := dao.Invoice.Ctx(backgroundCtx)
			if req.StartId != nil {
				query = query.WhereGTE(dao.Invoice.Columns().Id, req.StartId)
			} else if req.StartTime != nil {
				query = query.WhereGTE(dao.Invoice.Columns().CreateTime, req.StartTime)
			}
			if req.EndId != nil {
				query = query.WhereLTE(dao.Invoice.Columns().Id, req.EndId)
			} else if req.StartTime != nil {
				query = query.WhereLTE(dao.Invoice.Columns().CreateTime, req.EndTime)
			}
			query = query.WhereIn(dao.Invoice.Columns().IsDeleted, []int{0}).
				Limit(page*count, count).
				OmitEmpty()
			err = query.Scan(&list)

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

			// next page
			page = page + 1
			if list == nil || len(list) == 0 {
				break
			}
		}
	}()
	return &invoice.InternalWebhookSyncRes{}, nil
}
