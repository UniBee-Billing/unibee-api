package out

import (
	"fmt"
	"go-oversea-pay/internal/consts"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetPaymentWebhookEntranceUrl(pay *entity.OverseaPay) string {
	return fmt.Sprintf("%s/webhooks/%d/notifications?payId=%d", consts.GetConfigInstance().HostPath, pay.ChannelId, pay.Id)
}

func GetPaymentRedirectEntranceUrl(pay *entity.OverseaPay) string {
	return fmt.Sprintf("%s/redirect/%d/forward?payId=%d", consts.GetConfigInstance().HostPath, pay.ChannelId, pay.Id)
}
