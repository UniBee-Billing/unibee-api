package subscription

import "github.com/gogf/gf/v2/frame/g"

type TestClockWalkReq struct {
	g.Meta         `path:"/test_clock_walk" tags:"System-Admin" method:"post" summary:"Subscription Test CLock Walk"`
	SubscriptionId string `json:"subscriptionId" dc:"Subscription Id" v:"required"`
	NewTestClock   int64  `json:"newTestClock" dc:"NewTestClock" v:"required"`
}
type TestClockWalkRes struct {
}
