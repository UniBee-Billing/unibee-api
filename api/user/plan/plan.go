package plan

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/subscription/ro"
)

type SubscriptionPlanListReq struct {
	g.Meta     `path:"/subscription_plan_list" tags:"User-Plan-Controller" method:"post" summary:"订阅计划列表"`
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	Type       int    `p:"type"  dc:"不填查询所有类型，,1-main plan，2-addon plan" `
	Currency   string `p:"currency" d:"usd"  dc:"订阅计划货币"  `
	Page       int    `p:"page"  dc:"分页页码,0开始" `
	Count      int    `p:"count"  dc:"订阅计划货币" dc:"每页数量" `
}
type SubscriptionPlanListRes struct {
	Plans []*ro.PlanDetailRo `p:"plans" dc:"订阅计划明细"`
}
