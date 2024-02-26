package plan

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type SubscriptionPlanListReq struct {
	g.Meta   `path:"/subscription_plan_list" tags:"User-Plan-Controller" method:"post" summary:"Plan List"`
	Type     int    `p:"type"  dc:"Default All，,1-main plan，2-addon plan" `
	Currency string `p:"currency" dc:"Currency"  `
}
type SubscriptionPlanListRes struct {
	Plans []*ro.PlanDetailRo `p:"plans" dc:"Plan Detail"`
}
