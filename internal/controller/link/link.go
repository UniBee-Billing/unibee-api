package link

import (
	"fmt"
	"unibee/internal/cmd/config"
)

func GetInvoiceLink(invoiceId string) string {
	return fmt.Sprintf("%s/in/%s", config.GetConfigInstance().Server.GetServerPath(), invoiceId)
}

func GetPaymentLink(paymentId string) string {
	return fmt.Sprintf("%s/pay/%s", config.GetConfigInstance().Server.GetServerPath(), paymentId)
}
