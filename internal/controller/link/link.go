package link

import (
	"fmt"
	"unibee/internal/consts"
)

func GetInvoiceLink(invoiceId string) string {
	return fmt.Sprintf("%s/in/%s", consts.GetConfigInstance().Server.GetServerPath(), invoiceId)
}

func GetPaymentLink(paymentId string) string {
	return fmt.Sprintf("%s/pay/%s", consts.GetConfigInstance().Server.GetServerPath(), paymentId)
}
