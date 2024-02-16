package webhook

import "github.com/gogf/gf/v2/frame/g"

type SubscriptionWebhookCheckAndSetupReq struct {
	g.Meta `path:"/webhook_check_and_setup" tags:"Merchant-Setting-Controller" method:"post" summary:"Webhook Setup"`
}
type SubscriptionWebhookCheckAndSetupRes struct {
}
