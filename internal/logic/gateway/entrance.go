package gateway

import (
	"fmt"
	"net/url"
	"unibee/internal/cmd/config"
	entity "unibee/internal/model/entity/default"
)

func GetPaymentWebhookEntranceUrl(gatewayId uint64) string {
	return fmt.Sprintf("%s/payment/gateway_webhook_entry/%d/notifications", config.GetConfigInstance().Server.GetServerPath(), gatewayId)
}

func GetPaymentRedirectEntranceUrl(pay *entity.Payment) string {
	return fmt.Sprintf("%s/payment/redirect/%d/forward?paymentId=%s", config.GetConfigInstance().Server.GetServerPath(), pay.GatewayId, pay.PaymentId)
}

func GetPaymentMethodRedirectEntranceUrlCheckout(gatewayId uint64, success bool, subId string, redirectUrl string) string {
	return fmt.Sprintf("%s/payment/method/redirect/%d/forward?success=%v&session_id={CHECKOUT_SESSION_ID}&subId=%s&redirectUrl=%s", config.GetConfigInstance().Server.GetServerPath(), gatewayId, success, subId, url.QueryEscape(redirectUrl))
}

func GetPaymentRedirectEntranceUrlCheckout(pay *entity.Payment, success bool) string {
	invoiceId := ""
	if len(pay.InvoiceId) > 0 {
		invoiceId = pay.InvoiceId
	}
	if len(pay.SubscriptionId) > 0 {
		return fmt.Sprintf("%s/payment/redirect/%d/forward?paymentId=%s&subId=%s&invoiceId=%s&success=%v&session_id={CHECKOUT_SESSION_ID}", config.GetConfigInstance().Server.GetServerPath(), pay.GatewayId, pay.PaymentId, pay.SubscriptionId, invoiceId, success)
	} else {
		return fmt.Sprintf("%s/payment/redirect/%d/forward?paymentId=%s&invoiceId=%s&success=%v&session_id={CHECKOUT_SESSION_ID}", config.GetConfigInstance().Server.GetServerPath(), pay.GatewayId, pay.PaymentId, invoiceId, success)
	}
}

func GetPaypalPaymentRedirectEntranceUrlCheckout(pay *entity.Payment, success bool) string {
	invoiceId := ""
	if len(pay.InvoiceId) > 0 {
		invoiceId = pay.InvoiceId
	}
	if len(pay.SubscriptionId) > 0 {
		return fmt.Sprintf("%s/payment/redirect/%d/forward?paymentId=%s&subId=%s&invoiceId=%s&success=%v", config.GetConfigInstance().Server.GetServerPath(), pay.GatewayId, pay.PaymentId, pay.SubscriptionId, invoiceId, success)
	} else {
		return fmt.Sprintf("%s/payment/redirect/%d/forward?paymentId=%s&invoiceId=%s&success=%v", config.GetConfigInstance().Server.GetServerPath(), pay.GatewayId, pay.PaymentId, invoiceId, success)
	}
}
