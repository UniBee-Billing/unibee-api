package link

import (
	"context"
	"fmt"
	"unibee/internal/cmd/config"
	"unibee/internal/query"
	"unibee/utility"
)

func GetInvoiceLink(ctx context.Context, invoiceId string) string {
	if len(invoiceId) == 0 {
		return ""
	}
	return fmt.Sprintf("%s/in/%s?st=%s", config.GetConfigInstance().Server.GetServerPath(), invoiceId, GenerateInvoiceLinkSecurityToken(ctx, invoiceId))
}

func GetInvoicePdfLink(ctx context.Context, invoiceId string) string {
	if len(invoiceId) == 0 {
		return ""
	}
	return fmt.Sprintf("%s/in/pdf/%s?sm=%s", config.GetConfigInstance().Server.GetServerPath(), invoiceId, GenerateInvoiceLinkSecurityToken(ctx, invoiceId))
}

func GetPaymentLink(paymentId string) string {
	if len(paymentId) == 0 {
		return ""
	}
	return fmt.Sprintf("%s/pay/%s", config.GetConfigInstance().Server.GetServerPath(), paymentId)
}

func GenerateInvoiceLinkSecurityToken(ctx context.Context, invoiceId string) string {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	if one != nil {
		return utility.MD5(fmt.Sprintf("%d%s%s%d", one.CreateTime, one.UniqueId, one.InvoiceId, one.Id))
	} else {
		return invoiceId
	}
}

func VerifyInvoiceLinkSecurityToken(ctx context.Context, invoiceId string, token string) bool {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	if one == nil {
		return false
	}
	if token == invoiceId {
		return true
	}
	if token == utility.MD5(fmt.Sprintf("%d%s%s%d", one.CreateTime, one.UniqueId, one.InvoiceId, one.Id)) {
		return true
	}
	return false
}
