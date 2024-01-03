// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package webhook

import (
	"context"
	
	"go-oversea-pay/api/webhook/v1"
)

type IWebhookV1 interface {
	SubscriptionWebhookCheckAndSetup(ctx context.Context, req *v1.SubscriptionWebhookCheckAndSetupReq) (res *v1.SubscriptionWebhookCheckAndSetupRes, err error)
}


