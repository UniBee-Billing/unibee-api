package system

import (
	"context"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	invoice2 "unibee/internal/consumer/webhook/invoice"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/system/invoice"
)

func (c *ControllerInvoice) BatchSendInvoiceWebhookEvent(ctx context.Context, req *invoice.BatchSendInvoiceWebhookEventReq) (res *invoice.BatchSendInvoiceWebhookEventRes, err error) {
	utility.Assert(len(req.InvoiceIds) > 0, "Empty invoiceIds")

	for _, invoiceId := range req.InvoiceIds {
		one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
		if one != nil {
			if one.Status == consts.InvoiceStatusPending {
				one.Status = consts.InvoiceStatusPending
				invoice2.SendMerchantInvoiceWebhookBackground(one, event.UNIBEE_WEBHOOK_EVENT_INVOICE_CREATED, map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()})
			} else if one.Status == consts.InvoiceStatusProcessing {
				one.Status = consts.InvoiceStatusProcessing
				invoice2.SendMerchantInvoiceWebhookBackground(one, event.UNIBEE_WEBHOOK_EVENT_INVOICE_PROCESS, map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()})
			} else if one.Status == consts.InvoiceStatusPaid {
				one.Status = consts.InvoiceStatusPaid
				invoice2.SendMerchantInvoiceWebhookBackground(one, event.UNIBEE_WEBHOOK_EVENT_INVOICE_PAID, map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()})
			} else if one.Status == consts.InvoiceStatusCancelled {
				one.Status = consts.InvoiceStatusCancelled
				invoice2.SendMerchantInvoiceWebhookBackground(one, event.UNIBEE_WEBHOOK_EVENT_INVOICE_CANCELLED, map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()})
			} else if one.Status == consts.InvoiceStatusFailed {
				one.Status = consts.InvoiceStatusFailed
				invoice2.SendMerchantInvoiceWebhookBackground(one, event.UNIBEE_WEBHOOK_EVENT_INVOICE_FAILED, map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()})
			} else if one.Status == consts.InvoiceStatusReversed {
				one.Status = consts.InvoiceStatusReversed
				invoice2.SendMerchantInvoiceWebhookBackground(one, event.UNIBEE_WEBHOOK_EVENT_INVOICE_REVERSED, map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()})
			}
		}
	}

	return &invoice.BatchSendInvoiceWebhookEventRes{}, nil
}
