package out

import (
	"fmt"
	"go-oversea-pay/internal/consts"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetPaymentWebhookEntranceUrl(pay *entity.OverseaPay) string {
	return fmt.Sprintf("%s/webhooks/%s/notifications?payId=%d", consts.GetConfigInstance().HostPath, pay.ChannelPayId, pay.Id)
}

func GetPaymentRedirectEntranceUrl(pay *entity.OverseaPay) string {
	return fmt.Sprintf("%s/redirect/%s/forward?payId=%d", consts.GetConfigInstance().HostPath, pay.ChannelPayId, pay.Id)
}
