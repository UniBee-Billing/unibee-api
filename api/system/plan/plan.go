package plan

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type DetailReq struct {
	g.Meta `path:"/detail" tags:"System-Plan" method:"get,post" summary:"PlanDetail"`
	PlanId uint64 `json:"planId" dc:"PlanId" v:"required"`
}
type DetailRes struct {
	Plan *detail.PlanDetail `json:"plan" dc:"Plan Detail"`
}
