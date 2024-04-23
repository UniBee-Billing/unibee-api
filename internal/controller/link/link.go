package link

import (
	"context"
	"fmt"
	"unibee/internal/cmd/config"
)

func GetInvoiceLink(ctx context.Context, invoiceId string, st string) string {
	if len(invoiceId) == 0 {
		return ""
	}
	return fmt.Sprintf("%s/in/%s?st=%s", config.GetConfigInstance().Server.GetServerPath(), invoiceId, st)
}

func GetInvoicePdfLink(ctx context.Context, invoiceId string, st string) string {
	if len(invoiceId) == 0 {
		return ""
	}
	return fmt.Sprintf("%s/in/pdf/%s?sm=%s", config.GetConfigInstance().Server.GetServerPath(), invoiceId, st)
}

func GetPaymentLink(paymentId string) string {
	if len(paymentId) == 0 {
		return ""
	}
	return fmt.Sprintf("%s/pay/%s", config.GetConfigInstance().Server.GetServerPath(), paymentId)
}
