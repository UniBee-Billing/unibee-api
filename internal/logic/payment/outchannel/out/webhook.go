package out

import (
	"fmt"
	"go-oversea-pay/internal/consts"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetPaymentWebhookEntranceUrl(channelId int64) string {
	return fmt.Sprintf("%s/payment/webhooks/%d/notifications", consts.GetConfigInstance().HostPath, channelId)
}

//func GetPaymentWebhookEntranceUrlByPay(pay *entity.OverseaPay) string {
//	return fmt.Sprintf("%s/payment/webhooks/%d/notifications", consts.GetConfigInstance().HostPath, pay.ChannelId)
//}

func GetPaymentRedirectEntranceUrl(pay *entity.OverseaPay) string {
	return fmt.Sprintf("%s/payment/redirect/%d/forward?payId=%d", consts.GetConfigInstance().HostPath, pay.ChannelId, pay.Id)
}

func GetSubscriptionRedirectEntranceUrl(subscription *entity.Subscription, success bool) string {
	return fmt.Sprintf("%s/payment/redirect/%d/forward?subId=%d&success=%v&session_id={CHECKOUT_SESSION_ID}", consts.GetConfigInstance().HostPath, subscription.ChannelId, subscription.Id, success)
}
