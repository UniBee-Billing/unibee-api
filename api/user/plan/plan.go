package plan

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type ListReq struct {
	g.Meta   `path:"/list" tags:"User-Plan-Controller" method:"post" summary:"Plan List"`
	Type     int    `json:"type"  dc:"Default All，,1-main plan，2-addon plan" `
	Currency string `json:"currency" dc:"Currency"  `
}
type ListRes struct {
	Plans []*ro.PlanDetailRo `json:"plans" dc:"Plan Detail"`
}
