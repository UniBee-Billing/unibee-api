package out

import (
	"fmt"
	"go-oversea-pay/internal/consts"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetPaymentWebhookEntranceUrl(channelId int64) string {
	return fmt.Sprintf("%s/payment/channel_webhook_entry/%d/notifications", consts.GetConfigInstance().Server.DomainPath, channelId)
}

//func GetPaymentWebhookEntranceUrlByPay(pay *entity.OverseaPay) string {
//	return fmt.Sprintf("%s/payment/channel_webhook_entry/%d/notifications", consts.GetConfigInstance().HostPath, pay.ChannelId)
//}

func GetPaymentRedirectEntranceUrl(pay *entity.Payment) string {
	return fmt.Sprintf("%s/payment/redirect/%d/forward?payId=%d", consts.GetConfigInstance().Server.DomainPath, pay.ChannelId, pay.Id)
}

func GetSubscriptionRedirectEntranceUrl(subscription *entity.Subscription, success bool) string {
	return fmt.Sprintf("%s/payment/redirect/%d/forward?subId=%v&success=%v&session_id={CHECKOUT_SESSION_ID}", consts.GetConfigInstance().Server.DomainPath, subscription.ChannelId, subscription.SubscriptionId, success)
}
