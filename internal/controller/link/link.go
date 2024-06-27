package link

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/cmd/config"
)

func GetInvoiceLink(invoiceId string, st string) string {
	if len(invoiceId) == 0 {
		return ""
	}
	return fmt.Sprintf("%s/in/%s?st=%s&t=%d", config.GetConfigInstance().Server.GetServerPath(), invoiceId, st, gtime.Now().Timestamp())
}

func GetInvoicePdfLink(invoiceId string, st string) string {
	if len(invoiceId) == 0 {
		return ""
	}
	return fmt.Sprintf("%s/in/pdf/%s?st=%s&t=%d", config.GetConfigInstance().Server.GetServerPath(), invoiceId, st, gtime.Now().Timestamp())
}

func GetPaymentLink(paymentId string) string {
	if len(paymentId) == 0 {
		return ""
	}
	return fmt.Sprintf("%s/pay/%s", config.GetConfigInstance().Server.GetServerPath(), paymentId)
}

func GetExportLink(taskId int64) string {
	if taskId <= 0 {
		return ""
	}
	return fmt.Sprintf("%s/export/%v", config.GetConfigInstance().Server.GetServerPath(), taskId)
}
