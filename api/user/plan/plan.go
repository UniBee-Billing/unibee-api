package plan

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/gateway/ro"
)

type SubscriptionPlanListReq struct {
	g.Meta     `path:"/subscription_plan_list" tags:"User-Plan-Controller" method:"post" summary:"订阅计划列表"`
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	Type       int    `p:"type"  dc:"不填查询所有类型，,1-main plan，2-addon plan" `
	Currency   string `p:"currency" d:"usd"  dc:"订阅计划货币"  `
	SortField  string `p:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType   string `p:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page       int    `p:"page"  dc:"Page, Start WIth 0" `
	Count      int    `p:"count"  dc:"Count" dc:"Count Of Page" `
}
type SubscriptionPlanListRes struct {
	Plans []*ro.PlanDetailRo `p:"plans" dc:"订阅计划明细"`
}
