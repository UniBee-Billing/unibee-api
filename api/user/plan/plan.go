package plan

import (
	"github.com/gogf/gf/v2/frame/g"
	merhcnatPlan "unibee/api/bean/detail"
)

type ListReq struct {
	g.Meta    `path:"/list" tags:"User-Plan" method:"get,post" summary:"Plan List"`
	ProductId int64  `json:"productId"  dc:"filter id list of product, will use default product if not specified " `
	Type      []int  `json:"type" dc:"Default All，,1-main plan，2-addon plan" `
	Currency  string `json:"currency" dc:"Currency"  `
	Page      int    `json:"page"  dc:"Page, Start 0" `
	Count     int    `json:"count"  dc:"Count Of Per Page" `
}
type ListRes struct {
	Plans []*merhcnatPlan.PlanDetail `json:"plans" dc:"Plan Detail"`
	Total int                        `json:"total" dc:"Total"`
}
