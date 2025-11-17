package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type PreviewSubscriptionNextInvoiceReq struct {
	g.Meta         `path:"/preview_subscription_next_invoice" tags:"Subscription" method:"get" summary:"Subscription Next Invoice Preview"`
	UserId         uint64 `json:"userId" dc:"UserId"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId, unique, either ExternalUserId&Email or UserId needed if subscriptionId not specified"`
	ProductId      int64  `json:"productId" dc:"Id of product" dc:"default product will use if productId not specified and subscriptionId is blank"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId"`
}

type PreviewSubscriptionNextInvoiceRes struct {
	Subscription              *bean.Subscription                      `json:"subscription" dc:"Subscription"`
	Invoice                   *bean.Invoice                           `json:"invoice"`
	SubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"subscriptionPendingUpdate" dc:"subscriptionPendingUpdate"`
}

type ApplySubscriptionNextInvoiceReq struct {
	g.Meta                 `path:"/apply_subscription_next_invoice" tags:"Subscription" method:"post" summary:"Apply Discount Or Premo Credit To Next Invoice"`
	UserId                 uint64 `json:"userId" dc:"UserId"`
	ExternalUserId         string `json:"externalUserId" dc:"ExternalUserId, unique, either ExternalUserId&Email or UserId needed if subscriptionId not specified"`
	ProductId              int64  `json:"productId" dc:"Id of product" dc:"default product will use if productId not specified and subscriptionId is blank"`
	SubscriptionId         string `json:"subscriptionId" dc:"SubscriptionId"`
	DiscountCode           string `json:"discountCode" dc:"DiscountCode"`
	ApplyPromoCreditAmount int64  `json:"applyPromoCreditAmount"  dc:"apply promo credit amount"`
}

type ApplySubscriptionNextInvoiceRes struct {
	Subscription              *bean.Subscription                      `json:"subscription" dc:"Subscription"`
	Invoice                   *bean.Invoice                           `json:"invoice"`
	SubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"subscriptionPendingUpdate" dc:"subscriptionPendingUpdate"`
}
