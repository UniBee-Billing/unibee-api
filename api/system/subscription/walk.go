package subscription

import "github.com/gogf/gf/v2/frame/g"

type SubscriptionWalkTestClockReq struct {
	g.Meta         `path:"/subscription_test_clock_walk" tags:"System-Admin-Controller" method:"post" summary:"Subscription Test CLock Walk (In Process)"`
	SubscriptionId string `p:"subscriptionId" dc:"Subscription Id" v:"required"`
	NewWalkTime    int64  `p:"newWalkTime" dc:"NewWalkTime" v:"required"`
}
type SubscriptionWalkTestClockRes struct {
}