package gateway

import (
	"fmt"
	"unibee/internal/consts"
	entity "unibee/internal/model/entity/oversea_pay"
)

func GetPaymentWebhookEntranceUrl(gatewayId uint64) string {
	return fmt.Sprintf("%s/payment/gateway_webhook_entry/%d/notifications", consts.GetConfigInstance().Server.DomainPath, gatewayId)
}

//func GetPaymentWebhookEntranceUrlByPay(pay *entity.OverseaPay) string {
//	return fmt.Sprintf("%s/payment/gateway_webhook_entry/%d/notifications", consts.GetConfigInstance().HostPath, pay.Id)
//}

func GetPaymentRedirectEntranceUrl(pay *entity.Payment) string {
	return fmt.Sprintf("%s/payment/redirect/%d/forward?paymentId=%s", consts.GetConfigInstance().Server.DomainPath, pay.GatewayId, pay.PaymentId)
}

func GetPaymentRedirectEntranceUrlCheckout(pay *entity.Payment, success bool) string {
	if len(pay.SubscriptionId) > 0 {
		return fmt.Sprintf("%s/payment/redirect/%d/forward?paymentId=%s&subId=%s&success=%v&session_id={CHECKOUT_SESSION_ID}", consts.GetConfigInstance().Server.DomainPath, pay.GatewayId, pay.PaymentId, pay.SubscriptionId, success)
	} else {
		return fmt.Sprintf("%s/payment/redirect/%d/forward?paymentId=%s&success=%v&session_id={CHECKOUT_SESSION_ID}", consts.GetConfigInstance().Server.DomainPath, pay.GatewayId, pay.PaymentId, success)
	}
}

func GetSubscriptionRedirectEntranceUrl(subscription *entity.Subscription, success bool) string {
	return fmt.Sprintf("%s/payment/redirect/%d/forward?subId=%v&success=%v&session_id={CHECKOUT_SESSION_ID}", consts.GetConfigInstance().Server.DomainPath, subscription.GatewayId, subscription.SubscriptionId, success)
}
