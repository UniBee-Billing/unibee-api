package v1

import "github.com/gogf/gf/v2/frame/g"

type SubscriptionWebhookCheckAndSetupReq struct {
	g.Meta `path:"/subscription_webhook_check_and_setup" tags:"Subscription-Webhook-Admin-Controller" method:"post" summary:"Webhook 初始化"`
}
type SubscriptionWebhookCheckAndSetupRes struct {
}
