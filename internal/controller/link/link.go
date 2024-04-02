package link

import (
	"fmt"
	"unibee/internal/cmd/config"
)

func GetInvoiceLink(invoiceId string) string {
	if len(invoiceId) == 0 {
		return ""
	}
	return fmt.Sprintf("%s/in/%s", config.GetConfigInstance().Server.GetServerPath(), invoiceId)
}

func GetPaymentLink(paymentId string) string {
	if len(paymentId) == 0 {
		return ""
	}
	return fmt.Sprintf("%s/pay/%s", config.GetConfigInstance().Server.GetServerPath(), paymentId)
}
