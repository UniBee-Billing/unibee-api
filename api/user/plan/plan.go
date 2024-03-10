package plan

import (
	"github.com/gogf/gf/v2/frame/g"
	merhcnatPlan "unibee/api/merchant/plan"
)

type ListReq struct {
	g.Meta   `path:"/list" tags:"User-Plan" method:"get,post" summary:"Plan List"`
	Type     int    `json:"type"  dc:"Default All，,1-main plan，2-addon plan" `
	Currency string `json:"currency" dc:"Currency"  `
}
type ListRes struct {
	Plans []*merhcnatPlan.PlanDetail `json:"plans" dc:"Plan Detail"`
}
