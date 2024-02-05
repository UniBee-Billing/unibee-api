package plan

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/channel/ro"
)

type SubscriptionPlanListReq struct {
	g.Meta     `path:"/subscription_plan_list" tags:"User-Plan-Controller" method:"post" summary:"Plan List"`
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	Type       int    `p:"type"  dc:"Default All，,1-main plan，2-addon plan" `
	Currency   string `p:"currency" dc:"Currency"  `
	//SortField  string `p:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	//SortType   string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	//Page       int    `p:"page"  dc:"Page, Start WIth 0" `
	//Count      int    `p:"count"  dc:"Count" dc:"Count Of Page" `
}
type SubscriptionPlanListRes struct {
	Plans []*ro.PlanDetailRo `p:"plans" dc:"Plan Detail"`
}
