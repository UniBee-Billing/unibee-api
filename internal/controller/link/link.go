package link

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/cmd/config"
	entity "unibee/internal/model/entity/oversea_pay"
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

func GetTaskDownloadUrl(task *entity.MerchantBatchTask) string {
	if task == nil {
		return ""
	}
	return fmt.Sprintf("%s/export/%v", config.GetConfigInstance().Server.GetServerPath(), task.Id)
}
