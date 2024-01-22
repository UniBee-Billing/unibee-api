package subscription

import "github.com/gogf/gf/v2/frame/g"

type SubscriptionEndTrialReq struct {
	g.Meta         `path:"/subscription_end_trial" tags:"System-Admin-Controller" method:"post" summary:"Merchant 终止试用"`
	SubscriptionId string `p:"subscriptionId" dc:"订阅 ID" v:"required#请输入订阅 ID"`
}
type SubscriptionEndTrialRes struct {
}
