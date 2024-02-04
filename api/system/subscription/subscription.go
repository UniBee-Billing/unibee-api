package subscription

import "github.com/gogf/gf/v2/frame/g"

type SubscriptionEndTrialReq struct {
	g.Meta         `path:"/subscription_end_trial" tags:"System-Admin-Controller" method:"post" summary:"Merchant End Subscription Trial"`
	SubscriptionId string `p:"subscriptionId" dc:"Subscription Id" v:"required"`
}
type SubscriptionEndTrialRes struct {
}

type SubscriptionExpireReq struct {
	g.Meta         `path:"/subscription_expire" tags:"System-Admin-Controller" method:"post" summary:"Merchant Expire Subscription"`
	SubscriptionId string `p:"subscriptionId" dc:"Subscription Id" v:"required"`
}
type SubscriptionExpireRes struct {
}
