package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type PreviewSubscriptionNextInvoiceReq struct {
	g.Meta         `path:"/preview_subscription_next_invoice" tags:"Subscription" method:"get" summary:"Subscription Next Invoice Preview"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
}

type PreviewSubscriptionNextInvoiceRes struct {
	Subscription              *bean.Subscription                      `json:"subscription" dc:"Subscription"`
	Invoice                   *bean.Invoice                           `json:"invoice"`
	SubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"subscriptionPendingUpdate" dc:"subscriptionPendingUpdate"`
}
