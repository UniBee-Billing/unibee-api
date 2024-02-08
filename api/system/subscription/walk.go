package subscription

import "github.com/gogf/gf/v2/frame/g"

type SubscriptionWalkTestClockReq struct {
	g.Meta         `path:"/subscription_test_clock_walk" tags:"System-Admin-Controller" method:"post" summary:"Subscription Test CLock Walk (In Process)"`
	SubscriptionId string `p:"subscriptionId" dc:"Subscription Id" v:"required"`
	NewTestClock   int64  `p:"newTestClock" dc:"NewTestClock" v:"required"`
}
type SubscriptionWalkTestClockRes struct {
}
